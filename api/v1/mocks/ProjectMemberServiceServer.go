// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	mock "github.com/stretchr/testify/mock"
)

// ProjectMemberServiceServer is an autogenerated mock type for the ProjectMemberServiceServer type
type ProjectMemberServiceServer struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *ProjectMemberServiceServer) Create(_a0 context.Context, _a1 *v1.ProjectMemberCreateRequest) (*v1.ProjectMemberResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *v1.ProjectMemberResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberCreateRequest) (*v1.ProjectMemberResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberCreateRequest) *v1.ProjectMemberResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectMemberResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectMemberCreateRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *ProjectMemberServiceServer) Delete(_a0 context.Context, _a1 *v1.ProjectMemberDeleteRequest) (*v1.ProjectMemberResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 *v1.ProjectMemberResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberDeleteRequest) (*v1.ProjectMemberResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberDeleteRequest) *v1.ProjectMemberResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectMemberResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectMemberDeleteRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: _a0, _a1
func (_m *ProjectMemberServiceServer) Find(_a0 context.Context, _a1 *v1.ProjectMemberFindRequest) (*v1.ProjectMemberListResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 *v1.ProjectMemberListResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberFindRequest) (*v1.ProjectMemberListResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberFindRequest) *v1.ProjectMemberListResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectMemberListResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectMemberFindRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *ProjectMemberServiceServer) Get(_a0 context.Context, _a1 *v1.ProjectMemberGetRequest) (*v1.ProjectMemberResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *v1.ProjectMemberResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberGetRequest) (*v1.ProjectMemberResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberGetRequest) *v1.ProjectMemberResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectMemberResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectMemberGetRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *ProjectMemberServiceServer) Update(_a0 context.Context, _a1 *v1.ProjectMemberUpdateRequest) (*v1.ProjectMemberResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *v1.ProjectMemberResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberUpdateRequest) (*v1.ProjectMemberResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProjectMemberUpdateRequest) *v1.ProjectMemberResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProjectMemberResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProjectMemberUpdateRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProjectMemberServiceServer creates a new instance of ProjectMemberServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProjectMemberServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProjectMemberServiceServer {
	mock := &ProjectMemberServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}