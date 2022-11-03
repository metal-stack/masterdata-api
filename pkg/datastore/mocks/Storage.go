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
type Storage struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, ve
func (_m *Storage) Create(ctx context.Context, ve datastore.Entity) error {
	ret := _m.Called(ctx, ve)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, datastore.Entity) error); ok {
		r0 = rf(ctx, ve)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, ve
func (_m *Storage) Delete(ctx context.Context, ve datastore.Entity) error {
	ret := _m.Called(ctx, ve)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, datastore.Entity) error); ok {
		r0 = rf(ctx, ve)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: ctx, filter, paging, result
func (_m *Storage) Find(ctx context.Context, filter map[string]any, paging *v1.Paging, result any) (*uint64, error) {
	ret := _m.Called(ctx, filter, paging, result)

	var r0 *uint64
	if rf, ok := ret.Get(0).(func(context.Context, map[string]any, *v1.Paging, any) *uint64); ok {
		r0 = rf(ctx, filter, paging, result)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uint64)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, map[string]any, *v1.Paging, any) error); ok {
		r1 = rf(ctx, filter, paging, result)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, id, ve
func (_m *Storage) Get(ctx context.Context, id string, ve datastore.Entity) error {
	ret := _m.Called(ctx, id, ve)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, datastore.Entity) error); ok {
		r0 = rf(ctx, id, ve)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetHistory provides a mock function with given fields: ctx, id, at, ve
func (_m *Storage) GetHistory(ctx context.Context, id string, at time.Time, ve datastore.Entity) error {
	ret := _m.Called(ctx, id, at, ve)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, datastore.Entity) error); ok {
		r0 = rf(ctx, id, at, ve)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, ve
func (_m *Storage) Update(ctx context.Context, ve datastore.Entity) error {
	ret := _m.Called(ctx, ve)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, datastore.Entity) error); ok {
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
func NewStorage(t mockConstructorTestingTNewStorage) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
