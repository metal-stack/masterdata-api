package client

import (
	"context"

	"connectrpc.com/connect"
	compress "github.com/klauspost/connect-compress/v2"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/api/v1/apiv1connect"
)

// Client defines the client API
type (
	Client interface {
		Project() apiv1connect.ProjectServiceClient
		ProjectMember() apiv1connect.ProjectMemberServiceClient
		Tenant() apiv1connect.TenantServiceClient
		TenantMember() apiv1connect.TenantMemberServiceClient
		Version() apiv1connect.VersionServiceClient
	}

	client struct {
		config *Config
	}
)

// GRPCClient is a Client implementation with grpc transport.
func New(config *Config) Client {
	return &client{
		config: config,
	}
}

// Project is the root accessor for project related functions
func (c client) Project() apiv1connect.ProjectServiceClient {
	return apiv1connect.NewProjectServiceClient(
		c.config.HttpClient(),
		c.config.BaseURL,
		compress.WithAll(compress.LevelBalanced),
	)
}

// ProjectMember is the root accessor for project member related functions
func (c client) ProjectMember() apiv1connect.ProjectMemberServiceClient {
	return apiv1connect.NewProjectMemberServiceClient(
		c.config.HttpClient(),
		c.config.BaseURL,
		compress.WithAll(compress.LevelBalanced),
		connect.WithInterceptors(NamespaceInterceptor(c.config.Namespace)),
	)
}

// Tenant is the root accessor for tenant related functions
func (c client) Tenant() apiv1connect.TenantServiceClient {
	return apiv1connect.NewTenantServiceClient(
		c.config.HttpClient(),
		c.config.BaseURL,
		compress.WithAll(compress.LevelBalanced),
		connect.WithInterceptors(NamespaceInterceptor(c.config.Namespace)),
	)
}

// Tenant is the root accessor for tenant related functions
func (c client) TenantMember() apiv1connect.TenantMemberServiceClient {
	return apiv1connect.NewTenantMemberServiceClient(
		c.config.HttpClient(),
		c.config.BaseURL,
		compress.WithAll(compress.LevelBalanced),
		connect.WithInterceptors(NamespaceInterceptor(c.config.Namespace)),
	)
}

func (c client) Version() apiv1connect.VersionServiceClient {
	return apiv1connect.NewVersionServiceClient(
		c.config.HttpClient(),
		c.config.BaseURL,
		compress.WithAll(compress.LevelBalanced),
	)
}

func NamespaceInterceptor(namespace string) connect.UnaryInterceptorFunc {
	return func(uf connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
			switch r := ar.Any().(type) {
			case *v1.TenantMemberCreateRequest:
				if r.TenantMember.Namespace == "" {
					r.TenantMember.Namespace = namespace
				}
			case *v1.ProjectMemberCreateRequest:
				if r.ProjectMember.Namespace == "" {
					r.ProjectMember.Namespace = namespace
				}
			case *v1.TenantMemberFindRequest:
				if r.Namespace == "" {
					r.Namespace = namespace
				}
			case *v1.ProjectMemberFindRequest:
				if r.Namespace == "" {
					r.Namespace = namespace
				}
			case *v1.FindParticipatingProjectsRequest:
				if r.Namespace == "" {
					r.Namespace = namespace
				}
			case *v1.FindParticipatingTenantsRequest:
				if r.Namespace == "" {
					r.Namespace = namespace
				}
			case *v1.ListTenantMembersRequest:
				if r.Namespace == "" {
					r.Namespace = namespace
				}
			}
			return uf(ctx, ar)
		}
	}
}
