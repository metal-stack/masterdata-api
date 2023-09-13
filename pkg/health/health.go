package health

import (
	"context"

	"sync"

	"google.golang.org/grpc/codes"
	v1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// Server represents a Health Check server to check
// if a service is running or not on some host.
type Server struct {
	mu sync.Mutex
	// statusMap stores the serving status of the services this HealthServer monitors.
	statusMap map[string]v1.HealthCheckResponse_ServingStatus
}

// NewHealthServer creates a new health check server for grpc services.
func NewHealthServer() *Server {
	return &Server{
		statusMap: make(map[string]v1.HealthCheckResponse_ServingStatus),
	}
}

func (s *Server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	// no authentication required
	return ctx, nil
}

// Check checks if the grpc server is healthy and running.
func (s *Server) Check(ctx context.Context, in *v1.HealthCheckRequest) (*v1.HealthCheckResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.statusMap["initdb"]
	if !ok {
		return &v1.HealthCheckResponse{
			Status: v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
	_, ok = s.statusMap["migratedb"]
	if !ok {
		return &v1.HealthCheckResponse{
			Status: v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
	if in.Service == "" {
		// check the server overall health status.
		return &v1.HealthCheckResponse{
			Status: v1.HealthCheckResponse_SERVING,
		}, nil
	}
	if status, ok := s.statusMap[in.Service]; ok {
		return &v1.HealthCheckResponse{
			Status: status,
		}, nil
	}
	return nil, status.Errorf(codes.NotFound, "unknown service")
}

func (s *Server) Watch(*v1.HealthCheckRequest, v1.Health_WatchServer) error {
	// FIXME implement
	return nil
}

// SetServingStatus is called when need to reset the serving status of a service
// or insert a new service entry into the statusMap.
func (s *Server) SetServingStatus(service string, status v1.HealthCheckResponse_ServingStatus) {
	s.mu.Lock()
	s.statusMap[service] = status
	s.mu.Unlock()
}
