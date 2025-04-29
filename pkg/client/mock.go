package client

import (
	"github.com/metal-stack/masterdata-api/api/v1/apiv1connect"
)

type MockClient struct {
	psc  apiv1connect.ProjectServiceClient
	tsc  apiv1connect.TenantServiceClient
	pmsc apiv1connect.ProjectMemberServiceClient
	tmsc apiv1connect.TenantMemberServiceClient
	vsc  apiv1connect.VersionServiceClient
}

func NewMock(psc apiv1connect.ProjectServiceClient, tsc apiv1connect.TenantServiceClient, pmsc apiv1connect.ProjectMemberServiceClient, tmsc apiv1connect.TenantMemberServiceClient, vsc apiv1connect.VersionServiceClient) *MockClient {
	return &MockClient{
		psc:  psc,
		tsc:  tsc,
		pmsc: pmsc,
		tmsc: tmsc,
		vsc:  vsc,
	}
}

func (c *MockClient) Close() error {
	return nil
}
func (c *MockClient) Project() apiv1connect.ProjectServiceClient {
	return c.psc
}
func (c *MockClient) ProjectMember() apiv1connect.ProjectMemberServiceClient {
	return c.pmsc
}

func (c *MockClient) Tenant() apiv1connect.TenantServiceClient {
	return c.tsc
}
func (c *MockClient) TenantMember() apiv1connect.TenantMemberServiceClient {
	return c.tmsc
}
func (c *MockClient) Version() apiv1connect.VersionServiceClient {
	return c.vsc
}
