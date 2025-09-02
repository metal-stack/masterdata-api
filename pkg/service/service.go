package service

import (
	"context"
	"errors"
	"time"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/datastore"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	ProjectDataStore       datastore.Storage[*v1.Project]
	ProjectMemberDataStore datastore.Storage[*v1.ProjectMember]
	TenantDataStore        datastore.Storage[*v1.Tenant]
	TenantMemberDataStore  datastore.Storage[*v1.TenantMember]
)

type StorageStatusWrapper[E datastore.Entity] struct {
	storage datastore.Storage[E]
}

func NewStorageStatusWrapper[E datastore.Entity](s datastore.Storage[E]) datastore.Storage[E] {
	return StorageStatusWrapper[E]{
		storage: s,
	}
}

func (s StorageStatusWrapper[E]) Create(ctx context.Context, ve E) error {
	return wrapCreateStatusError(s.storage.Create(ctx, ve))
}

func (s StorageStatusWrapper[E]) Update(ctx context.Context, ve E) error {
	return wrapUpdateStatusError(s.storage.Update(ctx, ve))
}

func (s StorageStatusWrapper[E]) Delete(ctx context.Context, id string) error {
	return wrapDeleteStatusError(s.storage.Delete(ctx, id))
}

func (s StorageStatusWrapper[E]) DeleteAll(ctx context.Context, ids ...string) error {
	return wrapDeleteStatusError(s.storage.DeleteAll(ctx, ids...))
}

func (s StorageStatusWrapper[E]) Get(ctx context.Context, id string) (E, error) {
	e, err := s.storage.Get(ctx, id)
	return e, wrapGetStatusError(err)
}

func (s StorageStatusWrapper[E]) GetHistory(ctx context.Context, id string, at time.Time, ve E) error {
	return wrapGetStatusError(s.storage.GetHistory(ctx, id, at, ve))
}

func (s StorageStatusWrapper[E]) GetHistoryCreated(ctx context.Context, id string, ve E) error {
	return wrapGetStatusError(s.storage.GetHistoryCreated(ctx, id, ve))
}

func (s StorageStatusWrapper[E]) Find(ctx context.Context, paging *v1.Paging, filters ...any) ([]E, *uint64, error) {
	return s.storage.Find(ctx, paging, filters...)
}

// wrapCreateStatusError wraps some errors in a grpc status error
func wrapCreateStatusError(err error) error {
	if errors.As(err, &datastore.DuplicateKeyError{}) {
		err = status.Error(codes.AlreadyExists, err.Error())
	}
	return err
}

// wrapDeleteStatusError wraps some errors in a grpc status error
func wrapDeleteStatusError(err error) error {
	if errors.As(err, &datastore.NotFoundError{}) {
		err = status.Error(codes.NotFound, err.Error())
	} else if errors.As(err, &datastore.DataCorruptionError{}) {
		err = status.Error(codes.Internal, err.Error())
	}

	return err
}

// wrapGetStatusError wraps some errors in a grpc status error
func wrapGetStatusError(err error) error {
	if errors.As(err, &datastore.NotFoundError{}) {
		err = status.Error(codes.NotFound, err.Error())
	}
	return err
}

// wrapUpdateStatusError wraps some errors in a grpc status error
func wrapUpdateStatusError(err error) error {
	if errors.As(err, &datastore.OptimisticLockError{}) {
		err = status.Error(codes.FailedPrecondition, err.Error())
	}
	return err
}
