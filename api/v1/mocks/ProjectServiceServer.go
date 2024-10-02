// Code generated by mockery v2.46.1. DO NOT EDIT.

package mocks

import (
	context "context"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	mock "github.com/stretchr/testify/mock"
)

// ProjectServiceServer is an autogenerated mock type for the ProjectServiceServer type
type ProjectServiceServer struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *ProjectServiceServer) Create(_a0 context.Context, _a1 *v1.ProjectCreateRequest) (*v1.ProjectResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *v1.ProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectCreateRequest) (*v1.ProjectResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectCreateRequest) *v1.ProjectResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectCreateRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *ProjectServiceServer) Delete(_a0 context.Context, _a1 *v1.ProjectDeleteRequest) (*v1.ProjectResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 *v1.ProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectDeleteRequest) (*v1.ProjectResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectDeleteRequest) *v1.ProjectResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectDeleteRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: _a0, _a1
func (_m *ProjectServiceServer) Find(_a0 context.Context, _a1 *v1.ProjectFindRequest) (*v1.ProjectListResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 *v1.ProjectListResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectFindRequest) (*v1.ProjectListResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectFindRequest) *v1.ProjectListResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectListResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectFindRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *ProjectServiceServer) Get(_a0 context.Context, _a1 *v1.ProjectGetRequest) (*v1.ProjectResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *v1.ProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectGetRequest) (*v1.ProjectResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectGetRequest) *v1.ProjectResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectGetRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHistory provides a mock function with given fields: _a0, _a1
func (_m *ProjectServiceServer) GetHistory(_a0 context.Context, _a1 *v1.ProjectGetHistoryRequest) (*v1.ProjectResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetHistory")
	}

	var r0 *v1.ProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectGetHistoryRequest) (*v1.ProjectResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectGetHistoryRequest) *v1.ProjectResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectGetHistoryRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *ProjectServiceServer) Update(_a0 context.Context, _a1 *v1.ProjectUpdateRequest) (*v1.ProjectResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *v1.ProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectUpdateRequest) (*v1.ProjectResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectUpdateRequest) *v1.ProjectResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectUpdateRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProjectServiceServer creates a new instance of ProjectServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProjectServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProjectServiceServer {
	mock := &ProjectServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
