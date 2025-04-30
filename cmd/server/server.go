package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof" // nolint:gosec
	"time"

	"github.com/jmoiron/sqlx"
	apiv1 "github.com/metal-stack/masterdata-api/api/v1"

	apiv1connect "github.com/metal-stack/masterdata-api/api/v1/apiv1connect"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"github.com/metal-stack/masterdata-api/pkg/service"

	"github.com/metal-stack/v"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"

	"connectrpc.com/grpchealth"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type config struct {
	GrpcServerEndpoint     string
	MetricsEndpoint        string
	Log                    *slog.Logger
	ProjectDataStore       datastore.Storage[*apiv1.Project]
	ProjectMemberDataStore datastore.Storage[*apiv1.ProjectMember]
	TenantDataStore        datastore.Storage[*apiv1.Tenant]
	TenantMemberDataStore  datastore.Storage[*apiv1.TenantMember]
	DB                     *sqlx.DB
}
type server struct {
	c config

	projectDataStore       datastore.Storage[*apiv1.Project]
	projectMemberDataStore datastore.Storage[*apiv1.ProjectMember]
	tenantDataStore        datastore.Storage[*apiv1.Tenant]
	tenantMemberDataStore  datastore.Storage[*apiv1.TenantMember]
}

func newServer(c config) *server {
	return &server{
		c:                      c,
		projectDataStore:       c.ProjectDataStore,
		projectMemberDataStore: c.ProjectMemberDataStore,
		tenantDataStore:        c.TenantDataStore,
		tenantMemberDataStore:  c.TenantMemberDataStore,
	}
}
func (s *server) Run() error {
	s.c.Log.Info("starting masterdata-api", "version", v.V.String())

	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	exporter, err := prometheus.New()
	if err != nil {
		return err
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go func() {
		s.c.Log.Info("serving metrics", "at", fmt.Sprintf("%s/metrics", s.c.MetricsEndpoint))
		metricsServer := http.NewServeMux()
		metricsServer.Handle("/metrics", promhttp.Handler())
		ms := &http.Server{
			Addr:              s.c.MetricsEndpoint,
			Handler:           metricsServer,
			ReadHeaderTimeout: time.Minute,
		}
		err := ms.ListenAndServe()
		if err != nil {
			s.c.Log.Error("unable to start metric endpoint", "error", err)
			return
		}
	}()
	go func() {
		s.c.Log.Info("starting pprof endpoint of :2113")
		// inspect via
		// go tool pprof -http :8080 localhost:2113/debug/pprof/heap
		// go tool pprof -http :8080 localhost:2113/debug/pprof/goroutine
		server := http.Server{
			Addr:              ":2113",
			ReadHeaderTimeout: 1 * time.Minute,
		}
		err := server.ListenAndServe()
		if err != nil {
			s.c.Log.Error("failed to start pprof endpoint", "error", err)
			return
		}
	}()

	otelInterceptor, err := otelconnect.NewInterceptor(otelconnect.WithMeterProvider(provider))
	if err != nil {
		return err
	}

	loggingInterceptor := newLoggingInterceptor(s.c.Log)

	projectService := service.NewProjectService(s.c.Log, s.c.ProjectDataStore, s.c.ProjectMemberDataStore, s.c.TenantDataStore)
	projectMemberService := service.NewProjectMemberService(s.c.Log, s.c.ProjectDataStore, s.c.ProjectMemberDataStore, s.c.TenantDataStore)
	// FIXME db should not be required here
	tenantService := service.NewTenantService(s.c.DB, s.c.Log, s.c.TenantDataStore, s.c.TenantMemberDataStore)
	tenantMemberService := service.NewTenantMemberService(s.c.Log, s.c.TenantDataStore, s.c.TenantMemberDataStore)
	versionService := service.NewVersionService()

	// healthv1.RegisterHealthServer(grpcServer, healthServer)
	interceptors := connect.WithInterceptors(loggingInterceptor, otelInterceptor)

	mux := http.NewServeMux()
	mux.Handle(apiv1connect.NewProjectServiceHandler(projectService, interceptors))
	mux.Handle(apiv1connect.NewProjectMemberServiceHandler(projectMemberService, interceptors))
	mux.Handle(apiv1connect.NewTenantServiceHandler(tenantService, interceptors))
	mux.Handle(apiv1connect.NewTenantMemberServiceHandler(tenantMemberService, interceptors))
	mux.Handle(apiv1connect.NewVersionServiceHandler(versionService, interceptors))

	allServiceNames := []string{
		apiv1connect.ProjectServiceName,
		apiv1connect.ProjectMemberServiceName,
		apiv1connect.TenantServiceName,
		apiv1connect.TenantMemberServiceName,
		apiv1connect.VersionServiceName,
	}

	checker := grpchealth.NewStaticChecker(allServiceNames...)
	mux.Handle(grpchealth.NewHandler(checker))

	// enable remote service listing by enabling reflection
	reflector := grpcreflect.NewStaticReflector(allServiceNames...)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	server := http.Server{
		Addr: s.c.GrpcServerEndpoint,
		// For gRPC clients, it's convenient to support HTTP/2 without TLS. You can
		// avoid x/net/http2 by using http.ListenAndServeTLS.
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadHeaderTimeout: 1 * time.Minute,
	}

	s.c.Log.Info("started grpc server", "at", server.Addr)
	err = server.ListenAndServe()
	return err
}

func newLoggingInterceptor(log *slog.Logger) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			var (
				procedure = req.Spec().Procedure
				request   = req.Any()
			)
			if procedure == apiv1connect.VersionServiceGetProcedure {
				return next(ctx, req)
			}
			log.Debug("call", "proc", procedure, "req", request)

			response, err := next(ctx, req)
			if err != nil {
				log.Error("call", "proc", procedure, "error", err)
			} else {
				log.Debug("call", "proc", procedure, "req", request, "resp", response.Any())
			}

			return response, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
