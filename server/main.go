package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"os"
	"strings"
	"time"

	"google.golang.org/grpc/reflection"

	"github.com/metal-stack/masterdata-api/pkg/auth"
	"github.com/metal-stack/masterdata-api/pkg/health"
	"github.com/metal-stack/masterdata-api/pkg/interceptors/grpc_internalerror"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	apiv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"github.com/metal-stack/masterdata-api/pkg/service"
	"github.com/metal-stack/v"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	moduleName = "masterdata-api"
)

var (
	logger *slog.Logger
)

var rootCmd = &cobra.Command{
	Use:     moduleName,
	Short:   "api to manage masterdata data for metal cloud components",
	Version: v.V.String(),
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("failed executing root command", "error", err)
	}
}

func initConfig() {
	viper.SetEnvPrefix("MASTERDATA_API")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().IntP("port", "", 50051, "the port to serve on")
	rootCmd.Flags().BoolP("debug", "", false, "enable debugging")

	rootCmd.Flags().StringP("ca", "", "certs/ca.pem", "ca path")
	rootCmd.Flags().StringP("cert", "", "certs/server.pem", "server certificate path")
	rootCmd.Flags().StringP("certkey", "", "certs/server-key.pem", "server key path")

	rootCmd.Flags().StringP("dbhost", "", "localhost", "postgres database server hostname/ip")
	rootCmd.Flags().StringP("dbport", "", "5432", "postgres database server port")
	rootCmd.Flags().StringP("dbuser", "", "masterdata", "postgres database user")
	rootCmd.Flags().StringP("dbpassword", "", "password", "postgres database password")
	rootCmd.Flags().StringP("dbname", "", "masterdata", "postgres database name")
	rootCmd.Flags().StringP("dbsslmode", "", "disable", "sslmode to talk to the the database")
	rootCmd.Flags().StringP("hmackey", "", auth.HmacDefaultKey, "preshared hmac key to authenticate.")

	err := viper.BindPFlags(rootCmd.Flags())
	if err != nil {
		logger.Error("unable to construct root command", "error", err)
	}
}

func run() {

	lvl := slog.LevelInfo
	if viper.IsSet("debug") {
		lvl = slog.LevelDebug
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})

	logger := slog.New(jsonHandler)

	port := viper.GetInt("port")
	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("failed to listen", "error", err)
		panic(err)
	}

	logger.Info("starting masterdata-api", "version", v.V.String(), "address", addr)

	hmacKey := viper.GetString("hmackey")
	if hmacKey == "" {
		hmacKey = auth.HmacDefaultKey
	}
	auther, err := auth.NewHMACAuther(hmacKey, auth.EditUser)
	if err != nil {
		logger.Error("failed to create auther", "error", err)
		panic(err)
	}

	caFile := viper.GetString("ca")
	// Get system certificate pool
	certPool, err := x509.SystemCertPool()
	if err != nil {
		logger.Error("could not read system certificate pool", "error", err)
		panic(err)
	}

	if caFile != "" {
		logger.Info("using ca", "ca", caFile)
		ca, err := os.ReadFile(caFile)
		if err != nil {
			logger.Error("could not read ca certificate", "error", err)
			panic(err)
		}
		// Append the certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			logger.Error("failed to append ca certs", "error", err)
			panic(err)
		}
	}

	serverCert := viper.GetString("cert")
	serverKey := viper.GetString("certkey")
	cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		logger.Error("failed to load key pair", "error", err)
		panic(err)
	}

	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS12,
	})

	// FIXME migrate to grpc_middleware v2
	// grpcLogDeciderFunc := func(methodFullName string, err error) bool {
	// 	if err == nil && methodFullName == "/grpc.health.v1.Health/Check" {
	// 		return false
	// 	}
	// 	return true
	// }

	opts := []grpc.ServerOption{
		// Enable TLS for all incoming connections.
		grpc.Creds(creds),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			// FIXME migrate to grpc_middleware v2
			// grpc_zap.StreamServerInterceptor(logger, grpc_zap.WithDecider(grpcLogDeciderFunc)),
			grpc_auth.StreamServerInterceptor(auther.Auth),
			grpc_internalerror.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			// FIXME migrate to grpc_middleware v2
			//grpc_zap.UnaryServerInterceptor(logger, grpc_zap.WithDecider(grpcLogDeciderFunc)),
			grpc_auth.UnaryServerInterceptor(auther.Auth),
			grpc_internalerror.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	}

	// Set GRPC Interceptors
	// opts := []grpc.ServerOption{}
	// grpcServer := grpc.NewServer(opts...)
	grpcServer := grpc.NewServer(opts...)

	ves := []datastore.Entity{
		&apiv1.Project{},
		&apiv1.Tenant{},
	}
	dbHost := viper.GetString("dbhost")
	dbPort := viper.GetString("dbport")
	dbUser := viper.GetString("dbuser")
	dbPassword := viper.GetString("dbpassword")
	dbName := viper.GetString("dbname")
	dbSSLMode := viper.GetString("dbsslmode")

	db, err := datastore.NewPostgresDB(logger, dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode, ves...)
	if err != nil {
		logger.Error("failed to create postgres connection", "error", err)
		panic(err)
	}

	healthServer := health.NewHealthServer()

	err = datastore.Initdb(logger, db, healthServer, "initdb.d")
	if err != nil {
		logger.Error("unable to apply initdb content", "error", err)
	}

	err = datastore.MigrateDB(logger, db, healthServer)
	if err != nil {
		logger.Error("unable to apply migrate db", "error", err)
	}

	projectService, err := service.NewProjectService(db, logger)
	if err != nil {
		logger.Error("unable to create project service", "error", err)
		panic(err)
	}
	tenantService, err := service.NewTenantService(db, logger)
	if err != nil {
		logger.Error("unable to create tenant service", "error", err)
		panic(err)
	}

	apiv1.RegisterProjectServiceServer(grpcServer, projectService)
	apiv1.RegisterTenantServiceServer(grpcServer, tenantService)
	healthv1.RegisterHealthServer(grpcServer, healthServer)

	// After all your registrations, make sure all of the Prometheus metrics are initialized.
	grpc_prometheus.Register(grpcServer)
	// Register Prometheus metrics handler
	metricsServer := http.NewServeMux()
	metricsServer.Handle("/metrics", promhttp.Handler())
	go func() {
		logger.Info("starting metrics endpoint of :2112")
		server := http.Server{
			Addr:              ":2112",
			Handler:           metricsServer,
			ReadHeaderTimeout: 1 * time.Minute,
		}
		err := server.ListenAndServe()
		if err != nil {
			logger.Error("failed to start metrics endpoint", "error", err)
		}
		os.Exit(1)
	}()

	go func() {
		logger.Info("starting pprof endpoint of :2113")
		// inspect via
		// go tool pprof -http :8080 localhost:2113/debug/pprof/heap
		// go tool pprof -http :8080 localhost:2113/debug/pprof/goroutine
		server := http.Server{
			Addr:              ":2113",
			ReadHeaderTimeout: 1 * time.Minute,
		}
		err := server.ListenAndServe()
		if err != nil {
			logger.Error("failed to start pprof endpoint", "error", err)
		}
		os.Exit(1)
	}()

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("failed to serve", "error", err)
		panic(err)
	}
}
