package client

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"connectrpc.com/connect"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/api/v1/apiv1connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func Test_Client(t *testing.T) {
	const (
		namespace = "a"
	)

	var (
		// log                  = slog.Default()
		projectMemberService = &projectMemberService{}
		tenantMemberService  = &tenantMemberService{}

		interceptors = connect.WithInterceptors()
	)

	mux := http.NewServeMux()
	mux.Handle(apiv1connect.NewProjectMemberServiceHandler(projectMemberService, interceptors))
	mux.Handle(apiv1connect.NewTenantMemberServiceHandler(tenantMemberService, interceptors))

	lis, err := net.Listen("tcp", "")
	require.NoError(t, err)

	_, portString, err := net.SplitHostPort(lis.Addr().String())
	require.NoError(t, err)

	server := http.Server{
		Addr: lis.Addr().String(),
		// For gRPC clients, it's convenient to support HTTP/2 without TLS. You can
		// avoid x/net/http2 by using http.ListenAndServeTLS.
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadHeaderTimeout: 1 * time.Minute,
	}

	go func() {
		_ = server.Serve(lis)
	}()
	defer func() {
		err := server.Close()
		require.NoError(t, err)
	}()

	client := New(&Config{
		BaseURL:   fmt.Sprintf("http://localhost:%s", portString),
		Debug:     true,
		UserAgent: "sample-client",
		Namespace: namespace,
	})
	require.NoError(t, err)

	t.Run("check namespace interceptor sets missing namespace", func(t *testing.T) {
		t.Run("project member", func(t *testing.T) {
			projectMemberService.create = func(_ context.Context, pmcr *connect.Request[v1.ProjectMemberCreateRequest]) (*connect.Response[v1.ProjectMemberResponse], error) {
				assert.Equal(t, "project-a", pmcr.Msg.ProjectMember.ProjectId)
				assert.Equal(t, "tenant-a", pmcr.Msg.ProjectMember.TenantId)
				assert.Equal(t, namespace, pmcr.Msg.ProjectMember.Namespace)
				return connect.NewResponse(&v1.ProjectMemberResponse{}), nil
			}
			projectMemberService.find = func(_ context.Context, pmfr *connect.Request[v1.ProjectMemberFindRequest]) (*connect.Response[v1.ProjectMemberListResponse], error) {
				assert.Equal(t, namespace, pmfr.Msg.Namespace)
				return connect.NewResponse(&v1.ProjectMemberListResponse{}), nil

			}

			_, err = client.ProjectMember().Create(t.Context(), connect.NewRequest(&v1.ProjectMemberCreateRequest{
				ProjectMember: &v1.ProjectMember{
					ProjectId: "project-a",
					TenantId:  "tenant-a",
				},
			}))
			require.NoError(t, err)

			_, err = client.ProjectMember().Find(t.Context(), connect.NewRequest(&v1.ProjectMemberFindRequest{}))
			require.NoError(t, err)
		})

		t.Run("tenant member", func(t *testing.T) {
			tenantMemberService.create = func(_ context.Context, tmcr *connect.Request[v1.TenantMemberCreateRequest]) (*connect.Response[v1.TenantMemberResponse], error) {
				assert.Equal(t, "tenant-a", tmcr.Msg.TenantMember.TenantId)
				assert.Equal(t, namespace, tmcr.Msg.TenantMember.Namespace)
				return connect.NewResponse(&v1.TenantMemberResponse{}), nil
			}
			tenantMemberService.find = func(ctx context.Context, tmfr *connect.Request[v1.TenantMemberFindRequest]) (*connect.Response[v1.TenantMemberListResponse], error) {
				assert.Equal(t, namespace, tmfr.Msg.Namespace)
				return connect.NewResponse(&v1.TenantMemberListResponse{}), nil
			}

			_, err = client.TenantMember().Create(t.Context(), connect.NewRequest(&v1.TenantMemberCreateRequest{
				TenantMember: &v1.TenantMember{
					TenantId: "tenant-a",
				},
			}))
			require.NoError(t, err)

			_, err = client.TenantMember().Find(t.Context(), connect.NewRequest(&v1.TenantMemberFindRequest{}))
			require.NoError(t, err)
		})
	})

	t.Run("check explicit namespace can be set anyway", func(t *testing.T) {
		t.Run("project member", func(t *testing.T) {
			projectMemberService.create = func(_ context.Context, pmcr *connect.Request[v1.ProjectMemberCreateRequest]) (*connect.Response[v1.ProjectMemberResponse], error) {
				assert.Equal(t, "project-a", pmcr.Msg.ProjectMember.ProjectId)
				assert.Equal(t, "tenant-a", pmcr.Msg.ProjectMember.TenantId)
				assert.Equal(t, "b", pmcr.Msg.ProjectMember.Namespace)
				return connect.NewResponse(&v1.ProjectMemberResponse{}), nil
			}
			projectMemberService.find = func(_ context.Context, pmfr *connect.Request[v1.ProjectMemberFindRequest]) (*connect.Response[v1.ProjectMemberListResponse], error) {
				assert.Equal(t, "b", pmfr.Msg.Namespace)
				return connect.NewResponse(&v1.ProjectMemberListResponse{}), nil
			}

			_, err = client.ProjectMember().Create(t.Context(), connect.NewRequest(&v1.ProjectMemberCreateRequest{
				ProjectMember: &v1.ProjectMember{
					ProjectId: "project-a",
					TenantId:  "tenant-a",
					Namespace: "b",
				},
			}))
			require.NoError(t, err)

			_, err = client.ProjectMember().Find(t.Context(), connect.NewRequest(&v1.ProjectMemberFindRequest{
				Namespace: "b",
			}))
			require.NoError(t, err)
		})

		t.Run("tenant member", func(t *testing.T) {
			tenantMemberService.create = func(_ context.Context, tmcr *connect.Request[v1.TenantMemberCreateRequest]) (*connect.Response[v1.TenantMemberResponse], error) {
				assert.Equal(t, "tenant-a", tmcr.Msg.TenantMember.TenantId)
				assert.Equal(t, "b", tmcr.Msg.TenantMember.Namespace)
				return connect.NewResponse(&v1.TenantMemberResponse{}), nil
			}
			tenantMemberService.find = func(ctx context.Context, tmfr *connect.Request[v1.TenantMemberFindRequest]) (*connect.Response[v1.TenantMemberListResponse], error) {
				assert.Equal(t, "b", tmfr.Msg.Namespace)
				return connect.NewResponse(&v1.TenantMemberListResponse{}), nil
			}

			_, err = client.TenantMember().Create(t.Context(), connect.NewRequest(&v1.TenantMemberCreateRequest{
				TenantMember: &v1.TenantMember{
					TenantId:  "tenant-a",
					Namespace: "b",
				},
			}))
			require.NoError(t, err)

			_, err = client.TenantMember().Find(t.Context(), connect.NewRequest(&v1.TenantMemberFindRequest{
				Namespace: "b",
			}))
			require.NoError(t, err)
		})
	})
}

type projectMemberService struct {
	create func(context.Context, *connect.Request[v1.ProjectMemberCreateRequest]) (*connect.Response[v1.ProjectMemberResponse], error)
	find   func(context.Context, *connect.Request[v1.ProjectMemberFindRequest]) (*connect.Response[v1.ProjectMemberListResponse], error)
}

func (p *projectMemberService) Create(ctx context.Context, r *connect.Request[v1.ProjectMemberCreateRequest]) (*connect.Response[v1.ProjectMemberResponse], error) {
	return p.create(ctx, r)
}

func (p *projectMemberService) Delete(context.Context, *connect.Request[v1.ProjectMemberDeleteRequest]) (*connect.Response[v1.ProjectMemberResponse], error) {
	panic("unimplemented")
}

func (p *projectMemberService) Find(ctx context.Context, r *connect.Request[v1.ProjectMemberFindRequest]) (*connect.Response[v1.ProjectMemberListResponse], error) {
	return p.find(ctx, r)
}

func (p *projectMemberService) Get(context.Context, *connect.Request[v1.ProjectMemberGetRequest]) (*connect.Response[v1.ProjectMemberResponse], error) {
	panic("unimplemented")
}

func (p *projectMemberService) Update(context.Context, *connect.Request[v1.ProjectMemberUpdateRequest]) (*connect.Response[v1.ProjectMemberResponse], error) {
	panic("unimplemented")
}

type tenantMemberService struct {
	create func(context.Context, *connect.Request[v1.TenantMemberCreateRequest]) (*connect.Response[v1.TenantMemberResponse], error)
	find   func(context.Context, *connect.Request[v1.TenantMemberFindRequest]) (*connect.Response[v1.TenantMemberListResponse], error)
}

func (t *tenantMemberService) Create(ctx context.Context, r *connect.Request[v1.TenantMemberCreateRequest]) (*connect.Response[v1.TenantMemberResponse], error) {
	return t.create(ctx, r)
}

func (t *tenantMemberService) Delete(context.Context, *connect.Request[v1.TenantMemberDeleteRequest]) (*connect.Response[v1.TenantMemberResponse], error) {
	panic("unimplemented")
}

func (t *tenantMemberService) Find(ctx context.Context, r *connect.Request[v1.TenantMemberFindRequest]) (*connect.Response[v1.TenantMemberListResponse], error) {
	return t.find(ctx, r)
}

func (t *tenantMemberService) Get(context.Context, *connect.Request[v1.TenantMemberGetRequest]) (*connect.Response[v1.TenantMemberResponse], error) {
	panic("unimplemented")
}

func (t *tenantMemberService) Update(context.Context, *connect.Request[v1.TenantMemberUpdateRequest]) (*connect.Response[v1.TenantMemberResponse], error) {
	panic("unimplemented")
}
