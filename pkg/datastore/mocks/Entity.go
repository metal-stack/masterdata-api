// Code generated by mockery v2.50.4. DO NOT EDIT.

package mocks

import (
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	mock "github.com/stretchr/testify/mock"
)

// Entity is an autogenerated mock type for the Entity type
type Entity struct {
	mock.Mock
}

// APIVersion provides a mock function with no fields
func (_m *Entity) APIVersion() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for APIVersion")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetMeta provides a mock function with no fields
func (_m *Entity) GetMeta() *v1.Meta {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetMeta")
	}

	var r0 *v1.Meta
	if rf, ok := ret.Get(0).(func() *v1.Meta); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Meta)
		}
	}

	return r0
}

// JSONField provides a mock function with no fields
func (_m *Entity) JSONField() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for JSONField")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Kind provides a mock function with no fields
func (_m *Entity) Kind() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Kind")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Schema provides a mock function with no fields
func (_m *Entity) Schema() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Schema")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// TableName provides a mock function with no fields
func (_m *Entity) TableName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for TableName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewEntity creates a new instance of Entity. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEntity(t interface {
	mock.TestingT
	Cleanup(func())
}) *Entity {
	mock := &Entity{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
