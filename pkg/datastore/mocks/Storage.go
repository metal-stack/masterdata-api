// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	datastore "github.com/metal-stack/masterdata-api/pkg/datastore"
	mock "github.com/stretchr/testify/mock"

	time "time"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
)

// Storage is an autogenerated mock type for the Storage type
type Storage[E datastore.Entity] struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, ve
func (_m *Storage[E]) Create(ctx context.Context, ve E) error {
	ret := _m.Called(ctx, ve)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, E) error); ok {
		r0 = rf(ctx, ve)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *Storage[E]) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: ctx, filter, paging
func (_m *Storage[E]) Find(ctx context.Context, filter map[string]interface{}, paging *v1.Paging) ([]E, *uint64, error) {
	ret := _m.Called(ctx, filter, paging)

	var r0 []E
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, *v1.Paging) []E); ok {
		r0 = rf(ctx, filter, paging)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]E)
		}
	}

	var r1 *uint64
	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}, *v1.Paging) *uint64); ok {
		r1 = rf(ctx, filter, paging)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*uint64)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, map[string]interface{}, *v1.Paging) error); ok {
		r2 = rf(ctx, filter, paging)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Get provides a mock function with given fields: ctx, id
func (_m *Storage[E]) Get(ctx context.Context, id string) (E, error) {
	ret := _m.Called(ctx, id)

	var r0 E
	if rf, ok := ret.Get(0).(func(context.Context, string) E); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(E)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHistory provides a mock function with given fields: ctx, id, at, ve
func (_m *Storage[E]) GetHistory(ctx context.Context, id string, at time.Time, ve E) error {
	ret := _m.Called(ctx, id, at, ve)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, E) error); ok {
		r0 = rf(ctx, id, at, ve)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetHistoryCreated provides a mock function with given fields: ctx, id, ve
func (_m *Storage[E]) GetHistoryCreated(ctx context.Context, id string, ve E) error {
	ret := _m.Called(ctx, id, ve)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, E) error); ok {
		r0 = rf(ctx, id, ve)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, ve
func (_m *Storage[E]) Update(ctx context.Context, ve E) error {
	ret := _m.Called(ctx, ve)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, E) error); ok {
		r0 = rf(ctx, ve)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorage[E datastore.Entity](t mockConstructorTestingTNewStorage) *Storage[E] {
	mock := &Storage[E]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
