package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/metal-stack/masterdata-api/pkg/auth"
	"github.com/metal-stack/masterdata-api/pkg/health"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"

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
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl, AddSource: false})

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

	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	exemplarFromContext := func(ctx context.Context) prometheus.Labels {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return prometheus.Labels{"traceID": span.TraceID().String()}
		}
		return nil
	}
	allButHealthZ := func(ctx context.Context, callMeta interceptors.CallMeta) bool {
		return healthv1.Health_ServiceDesc.ServiceName != callMeta.Service
	}

	// Setup metric for panic recoveries.
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)
	reg.MustRegister(collectors.NewGoCollector())
	panicsTotal := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name: "grpc_req_panics_recovered_total",
		Help: "Total number of gRPC requests recovered from internal panic.",
	})
	grpcPanicRecoveryHandler := func(p any) (err error) {
		panicsTotal.Inc()
		logger.Error("msg", "recovered from panic", "panic", p, "stack", debug.Stack())
		return status.Errorf(codes.Internal, "%s", p)
	}

	// Set GRPC Interceptors
	// opts := []grpc.ServerOption{}
	// grpcServer := grpc.NewServer(opts...)
	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			// Order matters e.g. tracing interceptor have to create span first for the later exemplars to work.
			srvMetrics.UnaryServerInterceptor(),
			logging.UnaryServerInterceptor(interceptorLogger(logger)),
			selector.UnaryServerInterceptor(grpcauth.UnaryServerInterceptor(auther.Auth), selector.MatchFunc(allButHealthZ)),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			logging.StreamServerInterceptor(interceptorLogger(logger)),
			selector.StreamServerInterceptor(grpcauth.StreamServerInterceptor(auther.Auth), selector.MatchFunc(allButHealthZ)),
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	)

	ves := []datastore.Entity{
		&apiv1.Project{},
		&apiv1.ProjectMember{},
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
	projectMemberService, err := service.NewProjectMemberService(db, logger)
	if err != nil {
		logger.Error("unable to create project member service", "error", err)
		panic(err)

	}
	tenantService, err := service.NewTenantService(db, logger)
	if err != nil {
		logger.Error("unable to create tenant service", "error", err)
		panic(err)
	}

	apiv1.RegisterProjectServiceServer(grpcServer, projectService)
	apiv1.RegisterProjectMemberServiceServer(grpcServer, projectMemberService)
	apiv1.RegisterTenantServiceServer(grpcServer, tenantService)
	healthv1.RegisterHealthServer(grpcServer, healthServer)

	srvMetrics.InitializeMetrics(grpcServer)

	// Register Prometheus metrics handler
	metricsServer := http.NewServeMux()
	metricsServer.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
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

// interceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg, fields...)
		case logging.LevelInfo:
			l.Info(msg, fields...)
		case logging.LevelWarn:
			l.Warn(msg, fields...)
		case logging.LevelError:
			l.Error(msg, fields...)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
