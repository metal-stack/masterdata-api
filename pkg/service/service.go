package service

import (
	"context"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StorageStatusWrapper struct {
	storage datastore.Storage
}

func NewStorageStatusWrapper(s datastore.Storage) datastore.Storage {
	return StorageStatusWrapper{storage: s}
}

func (s StorageStatusWrapper) Create(ctx context.Context, ve datastore.VersionedJSONEntity) error {
	return wrapCreateStatusError(s.storage.Create(ctx, ve))
}

func (s StorageStatusWrapper) Update(ctx context.Context, ve datastore.VersionedJSONEntity) error {
	return wrapUpdateStatusError(s.storage.Update(ctx, ve))
}

func (s StorageStatusWrapper) Delete(ctx context.Context, ve datastore.VersionedJSONEntity) error {
	return wrapDeleteStatusError(s.storage.Delete(ctx, ve))
}

func (s StorageStatusWrapper) Get(ctx context.Context, id string, ve datastore.VersionedJSONEntity) error {
	return wrapGetStatusError(s.storage.Get(ctx, id, ve))
}

func (s StorageStatusWrapper) Find(ctx context.Context, filter map[string]interface{}, result interface{}) error {
	return s.storage.Find(ctx, filter, result)
}

// wrapCreateStatusError wraps some errors in a grpc status error
func wrapCreateStatusError(err error) error {
	switch err.(type) {
	case datastore.DuplicateKeyError:
		err = status.Error(codes.AlreadyExists, err.Error())
	}
	return err
}

// wrapDeleteStatusError wraps some errors in a grpc status error
func wrapDeleteStatusError(err error) error {
	switch err.(type) {
	case datastore.NotFoundError:
		err = status.Error(codes.NotFound, err.Error())
	case datastore.DataCorruptionError:
		err = status.Error(codes.Internal, err.Error())
	}
	return err
}

// wrapGetStatusError wraps some errors in a grpc status error
func wrapGetStatusError(err error) error {
	switch err.(type) {
	case datastore.NotFoundError:
		err = status.Error(codes.NotFound, err.Error())
	}
	return err
}

// wrapUpdateStatusError wraps some errors in a grpc status error
func wrapUpdateStatusError(err error) error {
	switch err.(type) {
	case datastore.OptimisticLockError:
		err = status.Error(codes.FailedPrecondition, err.Error())
	}
	return err
}
