// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
)

// ProjectServiceClient is an autogenerated mock type for the ProjectServiceClient type
type ProjectServiceClient struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, in, opts
func (_m *ProjectServiceClient) Create(ctx context.Context, in *v1.ProjectCreateRequest, opts ...grpc.CallOption) (*v1.ProjectResponse, error) {
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

	var r0 *v1.ProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectCreateRequest, ...grpc.CallOption) (*v1.ProjectResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectCreateRequest, ...grpc.CallOption) *v1.ProjectResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectCreateRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, in, opts
func (_m *ProjectServiceClient) Delete(ctx context.Context, in *v1.ProjectDeleteRequest, opts ...grpc.CallOption) (*v1.ProjectResponse, error) {
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

	var r0 *v1.ProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectDeleteRequest, ...grpc.CallOption) (*v1.ProjectResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectDeleteRequest, ...grpc.CallOption) *v1.ProjectResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectDeleteRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: ctx, in, opts
func (_m *ProjectServiceClient) Find(ctx context.Context, in *v1.ProjectFindRequest, opts ...grpc.CallOption) (*v1.ProjectListResponse, error) {
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

	var r0 *v1.ProjectListResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectFindRequest, ...grpc.CallOption) (*v1.ProjectListResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectFindRequest, ...grpc.CallOption) *v1.ProjectListResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectListResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectFindRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, in, opts
func (_m *ProjectServiceClient) Get(ctx context.Context, in *v1.ProjectGetRequest, opts ...grpc.CallOption) (*v1.ProjectResponse, error) {
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

	var r0 *v1.ProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectGetRequest, ...grpc.CallOption) (*v1.ProjectResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectGetRequest, ...grpc.CallOption) *v1.ProjectResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectGetRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHistory provides a mock function with given fields: ctx, in, opts
func (_m *ProjectServiceClient) GetHistory(ctx context.Context, in *v1.ProjectGetHistoryRequest, opts ...grpc.CallOption) (*v1.ProjectResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetHistory")
	}

	var r0 *v1.ProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectGetHistoryRequest, ...grpc.CallOption) (*v1.ProjectResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectGetHistoryRequest, ...grpc.CallOption) *v1.ProjectResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectGetHistoryRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, in, opts
func (_m *ProjectServiceClient) Update(ctx context.Context, in *v1.ProjectUpdateRequest, opts ...grpc.CallOption) (*v1.ProjectResponse, error) {
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

	var r0 *v1.ProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectUpdateRequest, ...grpc.CallOption) (*v1.ProjectResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectUpdateRequest, ...grpc.CallOption) *v1.ProjectResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectUpdateRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProjectServiceClient creates a new instance of ProjectServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProjectServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProjectServiceClient {
	mock := &ProjectServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
