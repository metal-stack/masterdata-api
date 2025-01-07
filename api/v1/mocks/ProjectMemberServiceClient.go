// Code generated by mockery v2.50.4. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
)

// ProjectMemberServiceClient is an autogenerated mock type for the ProjectMemberServiceClient type
type ProjectMemberServiceClient struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, in, opts
func (_m *ProjectMemberServiceClient) Create(ctx context.Context, in *v1.ProjectMemberCreateRequest, opts ...grpc.CallOption) (*v1.ProjectMemberResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *v1.ProjectMemberResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberCreateRequest, ...grpc.CallOption) (*v1.ProjectMemberResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberCreateRequest, ...grpc.CallOption) *v1.ProjectMemberResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectMemberResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectMemberCreateRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, in, opts
func (_m *ProjectMemberServiceClient) Delete(ctx context.Context, in *v1.ProjectMemberDeleteRequest, opts ...grpc.CallOption) (*v1.ProjectMemberResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 *v1.ProjectMemberResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberDeleteRequest, ...grpc.CallOption) (*v1.ProjectMemberResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberDeleteRequest, ...grpc.CallOption) *v1.ProjectMemberResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectMemberResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectMemberDeleteRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: ctx, in, opts
func (_m *ProjectMemberServiceClient) Find(ctx context.Context, in *v1.ProjectMemberFindRequest, opts ...grpc.CallOption) (*v1.ProjectMemberListResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 *v1.ProjectMemberListResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberFindRequest, ...grpc.CallOption) (*v1.ProjectMemberListResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberFindRequest, ...grpc.CallOption) *v1.ProjectMemberListResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectMemberListResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectMemberFindRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, in, opts
func (_m *ProjectMemberServiceClient) Get(ctx context.Context, in *v1.ProjectMemberGetRequest, opts ...grpc.CallOption) (*v1.ProjectMemberResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *v1.ProjectMemberResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberGetRequest, ...grpc.CallOption) (*v1.ProjectMemberResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberGetRequest, ...grpc.CallOption) *v1.ProjectMemberResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectMemberResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectMemberGetRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, in, opts
func (_m *ProjectMemberServiceClient) Update(ctx context.Context, in *v1.ProjectMemberUpdateRequest, opts ...grpc.CallOption) (*v1.ProjectMemberResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *v1.ProjectMemberResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberUpdateRequest, ...grpc.CallOption) (*v1.ProjectMemberResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberUpdateRequest, ...grpc.CallOption) *v1.ProjectMemberResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectMemberResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectMemberUpdateRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProjectMemberServiceClient creates a new instance of ProjectMemberServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProjectMemberServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProjectMemberServiceClient {
	mock := &ProjectMemberServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
