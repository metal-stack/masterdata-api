// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UnsafeTenantServiceServer is an autogenerated mock type for the UnsafeTenantServiceServer type
type UnsafeTenantServiceServer struct {
	mock.Mock
}

// mustEmbedUnimplementedTenantServiceServer provides a mock function with given fields:
func (_m *UnsafeTenantServiceServer) mustEmbedUnimplementedTenantServiceServer() {
	_m.Called()
}

// NewUnsafeTenantServiceServer creates a new instance of UnsafeTenantServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUnsafeTenantServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *UnsafeTenantServiceServer {
	mock := &UnsafeTenantServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}