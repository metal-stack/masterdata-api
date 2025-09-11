package client

import (
	"context"
	"log/slog"
	"net"
	"strconv"
	"testing"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Test_Client(t *testing.T) {
	const (
		namespace = "a"
	)

	var (
		log                 = slog.Default()
		grpcServer          = grpc.NewServer()
		projectMemberServer = &projectMemberServer{}
		tenantMemberServer  = &tenantMemberServer{}
	)

	v1.RegisterProjectMemberServiceServer(grpcServer, projectMemberServer)
	v1.RegisterTenantMemberServiceServer(grpcServer, tenantMemberServer)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", "")
	require.NoError(t, err)

	go func() {
		err = grpcServer.Serve(lis)
		require.NoError(t, err)
	}()
	defer func() {
		grpcServer.Stop()
	}()

	_, portString, err := net.SplitHostPort(lis.Addr().String())
	require.NoError(t, err)

	port, err := strconv.Atoi(portString)
	require.NoError(t, err)

	client, err := NewClient(&Config{
		Hostname:  "localhost",
		Port:      uint(port),
		Insecure:  true,
		Logger:    log,
		Namespace: namespace,
	})
	require.NoError(t, err)

	t.Run("check namespace interceptor sets missing namespace", func(t *testing.T) {
		t.Run("project member", func(t *testing.T) {
			projectMemberServer.create = func(ctx context.Context, pmcr *v1.ProjectMemberCreateRequest) (*v1.ProjectMemberResponse, error) {
				assert.Equal(t, "project-a", pmcr.ProjectMember.ProjectId)
				assert.Equal(t, "tenant-a", pmcr.ProjectMember.TenantId)
				assert.Equal(t, namespace, pmcr.ProjectMember.Namespace)
				return &v1.ProjectMemberResponse{}, nil
			}
			projectMemberServer.find = func(ctx context.Context, pmfr *v1.ProjectMemberFindRequest) (*v1.ProjectMemberListResponse, error) {
				assert.Equal(t, namespace, pmfr.Namespace)
				return &v1.ProjectMemberListResponse{}, nil
			}

			_, err = client.ProjectMember().Create(t.Context(), &v1.ProjectMemberCreateRequest{
				ProjectMember: &v1.ProjectMember{
					ProjectId: "project-a",
					TenantId:  "tenant-a",
				},
			})
			require.NoError(t, err)

			_, err = client.ProjectMember().Find(t.Context(), &v1.ProjectMemberFindRequest{})
			require.NoError(t, err)
		})

		t.Run("tenant member", func(t *testing.T) {
			tenantMemberServer.create = func(ctx context.Context, tmcr *v1.TenantMemberCreateRequest) (*v1.TenantMemberResponse, error) {
				assert.Equal(t, "tenant-a", tmcr.TenantMember.TenantId)
				assert.Equal(t, namespace, tmcr.TenantMember.Namespace)
				return &v1.TenantMemberResponse{}, nil
			}
			tenantMemberServer.find = func(ctx context.Context, tmfr *v1.TenantMemberFindRequest) (*v1.TenantMemberListResponse, error) {
				assert.Equal(t, namespace, tmfr.Namespace)
				return &v1.TenantMemberListResponse{}, nil
			}

			_, err = client.TenantMember().Create(t.Context(), &v1.TenantMemberCreateRequest{
				TenantMember: &v1.TenantMember{
					TenantId: "tenant-a",
				},
			})
			require.NoError(t, err)

			_, err = client.TenantMember().Find(t.Context(), &v1.TenantMemberFindRequest{})
			require.NoError(t, err)
		})
	})

	t.Run("check explicit namespace can be set anyway", func(t *testing.T) {
		t.Run("project member", func(t *testing.T) {
			projectMemberServer.create = func(ctx context.Context, pmcr *v1.ProjectMemberCreateRequest) (*v1.ProjectMemberResponse, error) {
				assert.Equal(t, "project-a", pmcr.ProjectMember.ProjectId)
				assert.Equal(t, "tenant-a", pmcr.ProjectMember.TenantId)
				assert.Equal(t, "b", pmcr.ProjectMember.Namespace)
				return &v1.ProjectMemberResponse{}, nil
			}
			projectMemberServer.find = func(ctx context.Context, pmfr *v1.ProjectMemberFindRequest) (*v1.ProjectMemberListResponse, error) {
				assert.Equal(t, "b", pmfr.Namespace)
				return &v1.ProjectMemberListResponse{}, nil
			}

			_, err = client.ProjectMember().Create(t.Context(), &v1.ProjectMemberCreateRequest{
				ProjectMember: &v1.ProjectMember{
					ProjectId: "project-a",
					TenantId:  "tenant-a",
					Namespace: "b",
				},
			})
			require.NoError(t, err)

			_, err = client.ProjectMember().Find(t.Context(), &v1.ProjectMemberFindRequest{
				Namespace: "b",
			})
			require.NoError(t, err)
		})

		t.Run("tenant member", func(t *testing.T) {
			tenantMemberServer.create = func(ctx context.Context, tmcr *v1.TenantMemberCreateRequest) (*v1.TenantMemberResponse, error) {
				assert.Equal(t, "tenant-a", tmcr.TenantMember.TenantId)
				assert.Equal(t, "b", tmcr.TenantMember.Namespace)
				return &v1.TenantMemberResponse{}, nil
			}
			tenantMemberServer.find = func(ctx context.Context, tmfr *v1.TenantMemberFindRequest) (*v1.TenantMemberListResponse, error) {
				assert.Equal(t, "b", tmfr.Namespace)
				return &v1.TenantMemberListResponse{}, nil
			}

			_, err = client.TenantMember().Create(t.Context(), &v1.TenantMemberCreateRequest{
				TenantMember: &v1.TenantMember{
					TenantId:  "tenant-a",
					Namespace: "b",
				},
			})
			require.NoError(t, err)

			_, err = client.TenantMember().Find(t.Context(), &v1.TenantMemberFindRequest{
				Namespace: "b",
			})
			require.NoError(t, err)
		})
	})
}

type projectMemberServer struct {
	create func(context.Context, *v1.ProjectMemberCreateRequest) (*v1.ProjectMemberResponse, error)
	find   func(context.Context, *v1.ProjectMemberFindRequest) (*v1.ProjectMemberListResponse, error)
}

func (t *projectMemberServer) Create(ctx context.Context, r *v1.ProjectMemberCreateRequest) (*v1.ProjectMemberResponse, error) {
	return t.create(ctx, r)
}

func (t *projectMemberServer) Delete(context.Context, *v1.ProjectMemberDeleteRequest) (*v1.ProjectMemberResponse, error) {
	panic("unimplemented")
}

func (t *projectMemberServer) Find(ctx context.Context, r *v1.ProjectMemberFindRequest) (*v1.ProjectMemberListResponse, error) {
	return t.find(ctx, r)
}

func (t *projectMemberServer) Get(context.Context, *v1.ProjectMemberGetRequest) (*v1.ProjectMemberResponse, error) {
	panic("unimplemented")
}

func (t *projectMemberServer) Update(context.Context, *v1.ProjectMemberUpdateRequest) (*v1.ProjectMemberResponse, error) {
	panic("unimplemented")
}

type tenantMemberServer struct {
	create func(context.Context, *v1.TenantMemberCreateRequest) (*v1.TenantMemberResponse, error)
	find   func(context.Context, *v1.TenantMemberFindRequest) (*v1.TenantMemberListResponse, error)
}

func (t *tenantMemberServer) Create(ctx context.Context, r *v1.TenantMemberCreateRequest) (*v1.TenantMemberResponse, error) {
	return t.create(ctx, r)
}

func (t *tenantMemberServer) Delete(context.Context, *v1.TenantMemberDeleteRequest) (*v1.TenantMemberResponse, error) {
	panic("unimplemented")
}

func (t *tenantMemberServer) Find(ctx context.Context, r *v1.TenantMemberFindRequest) (*v1.TenantMemberListResponse, error) {
	return t.find(ctx, r)
}

func (t *tenantMemberServer) Get(context.Context, *v1.TenantMemberGetRequest) (*v1.TenantMemberResponse, error) {
	panic("unimplemented")
}

func (t *tenantMemberServer) Update(context.Context, *v1.TenantMemberUpdateRequest) (*v1.TenantMemberResponse, error) {
	panic("unimplemented")
}
