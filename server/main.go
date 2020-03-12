package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"google.golang.org/grpc/reflection"

	"github.com/metal-stack/masterdata-api/pkg/auth"
	"github.com/metal-stack/masterdata-api/pkg/health"
	"github.com/metal-stack/masterdata-api/pkg/interceptors/grpc_internalerror"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	healthv1 "github.com/metal-stack/masterdata-api/api/grpc/health/v1"
	apiv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"github.com/metal-stack/masterdata-api/pkg/service"
	"github.com/metal-stack/v"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	moduleName = "masterdata-api"
)

var (
	logger *zap.Logger
)

var rootCmd = &cobra.Command{
	Use:     moduleName,
	Short:   "an api manage masterdata data for metal cloud components",
	Version: v.V.String(),
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Sugar().Errorw("failed executing root command", "error", err)
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
		logger.Sugar().Errorw("unable to construct root command:%v", err)
	}
}

func run() {
	logger, _ = zap.NewProduction()
	sLogger := logger.Sugar()
	defer func() {
		err := logger.Sync() // flushes buffer, if any
		if err != nil {
			fmt.Printf("unable to sync logger buffers:%v", err)
		}
	}()

	port := viper.GetInt("port")
	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	sLogger.Infow("starting masterdata-api", "version", v.V.String(), "address", addr)

	hmacKey := viper.GetString("hmackey")
	if hmacKey == "" {
		hmacKey = auth.HmacDefaultKey
	}
	auther, err := auth.NewHMACAuther(logger, hmacKey, auth.EditUser)
	if err != nil {
		logger.Fatal("failed to create auther", zap.Error(err))
	}

	caFile := viper.GetString("ca")
	// Get system certificate pool
	certPool, err := x509.SystemCertPool()
	if err != nil {
		logger.Fatal("could not read system certificate pool", zap.Error(err))
	}

	if caFile != "" {
		sLogger.Info("using ca", "ca", caFile)
		ca, err := ioutil.ReadFile(caFile)
		if err != nil {
			logger.Fatal("could not read ca certificate", zap.Error(err))
		}
		// Append the certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			logger.Fatal("failed to append ca certs", zap.Error(err))
		}
	}

	serverCert := viper.GetString("cert")
	serverKey := viper.GetString("certkey")
	cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		logger.Fatal("failed to load key pair", zap.Error(err))
	}

	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
	})

	opts := []grpc.ServerOption{
		// Enable TLS for all incoming connections.
		grpc.Creds(creds),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(logger),
			grpc_auth.StreamServerInterceptor(auther.Auth),
			grpc_internalerror.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_auth.UnaryServerInterceptor(auther.Auth),
			grpc_internalerror.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	}

	// Set GRPC Interceptors
	// opts := []grpc.ServerOption{}
	// grpcServer := grpc.NewServer(opts...)
	grpcServer := grpc.NewServer(opts...)

	ves := []datastore.VersionedJSONEntity{
		&apiv1.Project{},
		&apiv1.Tenant{},
	}
	dbHost := viper.GetString("dbhost")
	dbPort := viper.GetString("dbport")
	dbUser := viper.GetString("dbuser")
	dbPassword := viper.GetString("dbpassword")
	dbName := viper.GetString("dbname")
	dbSSLMode := viper.GetString("dbsslmode")

	storage, err := datastore.NewPostgresStorage(logger, dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode, ves...)
	if err != nil {
		logger.Fatal("failed to create postgres connection", zap.Error(err))
	}

	healthServer := health.NewHealthServer()

	err = storage.Initdb(healthServer, "initdb.d")
	if err != nil {
		sLogger.Errorw("unable to apply initdb content", "error", err)
	}

	projectService := service.NewProjectService(storage, logger)
	tenantService := service.NewTenantService(storage, logger)

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
		err := http.ListenAndServe(":2112", metricsServer)
		if err != nil {
			sLogger.Errorw("failed to start metrics endpoint", "error", err)
		}
		os.Exit(1)
	}()

	go func() {
		logger.Info("starting pprof endpoint of :2113")
		// inspect via
		// go tool pprof -http :8080 localhost:2113/debug/pprof/heap
		// go tool pprof -http :8080 localhost:2113/debug/pprof/goroutine
		err := http.ListenAndServe(":2113", nil)
		if err != nil {
			sLogger.Errorw("failed to start pprof endpoint", "error", err)
		}
		os.Exit(1)
	}()

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
