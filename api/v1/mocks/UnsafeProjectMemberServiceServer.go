// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UnsafeProjectMemberServiceServer is an autogenerated mock type for the UnsafeProjectMemberServiceServer type
type UnsafeProjectMemberServiceServer struct {
	mock.Mock
}

// mustEmbedUnimplementedProjectMemberServiceServer provides a mock function with no fields
func (_m *UnsafeProjectMemberServiceServer) mustEmbedUnimplementedProjectMemberServiceServer() {
	_m.Called()
}

// NewUnsafeProjectMemberServiceServer creates a new instance of UnsafeProjectMemberServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUnsafeProjectMemberServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *UnsafeProjectMemberServiceServer {
	mock := &UnsafeProjectMemberServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
