package client

import (
	v1 "github.com/metal-stack/masterdata-api/api/v1"
)

type MockClient struct {
	psc  v1.ProjectServiceClient
	tsc  v1.TenantServiceClient
	pmsc v1.ProjectMemberServiceClient
	tmsc v1.TenantMemberServiceClient
}

func NewMock(psc v1.ProjectServiceClient, tsc v1.TenantServiceClient, pmsc v1.ProjectMemberServiceClient, tmsc v1.TenantMemberServiceClient) *MockClient {
	return &MockClient{
		psc:  psc,
		tsc:  tsc,
		pmsc: pmsc,
		tmsc: tmsc,
	}
}

func (c *MockClient) Close() error {
	return nil
}
func (c *MockClient) Project() v1.ProjectServiceClient {
	return c.psc
}
func (c *MockClient) ProjectMember() v1.ProjectMemberServiceClient {
	return c.pmsc
}

func (c *MockClient) Tenant() v1.TenantServiceClient {
	return c.tsc
}
func (c *MockClient) TenantMember() v1.TenantMemberServiceClient {
	return c.tmsc
}
