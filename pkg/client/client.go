package client

import (
	compress "github.com/klauspost/connect-compress/v2"
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
		config DialConfig
	}
)

// GRPCClient is a Client implementation with grpc transport.
func New(config DialConfig) Client {
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
	)
}

// Tenant is the root accessor for tenant related functions
func (c client) Tenant() apiv1connect.TenantServiceClient {
	return apiv1connect.NewTenantServiceClient(
		c.config.HttpClient(),
		c.config.BaseURL,
		compress.WithAll(compress.LevelBalanced),
	)
}

// Tenant is the root accessor for tenant related functions
func (c client) TenantMember() apiv1connect.TenantMemberServiceClient {
	return apiv1connect.NewTenantMemberServiceClient(
		c.config.HttpClient(),
		c.config.BaseURL,
		compress.WithAll(compress.LevelBalanced),
	)
}

func (c client) Version() apiv1connect.VersionServiceClient {
	return apiv1connect.NewVersionServiceClient(
		c.config.HttpClient(),
		c.config.BaseURL,
		compress.WithAll(compress.LevelBalanced),
	)
}
