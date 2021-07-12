package datastore

import (
	"errors"

	"github.com/lib/pq"
)

// OptimisticLockError indicates that the operation could not be executed because the dataset to update has changed in the meantime.
// clients can decide to read the current dataset and retry the operation.
type OptimisticLockError struct {
	msg string
}

func (o OptimisticLockError) Error() string {
	return o.msg
}

// NewOptimisticLockError is called to create an OptimisticLockError error
func NewOptimisticLockError(msg string) OptimisticLockError {
	return OptimisticLockError{msg: msg}
}

// DuplicateKeyError indicates that an entity with the given id already exists
type DuplicateKeyError struct {
	msg string
}

func (o DuplicateKeyError) Error() string {
	return o.msg
}

// NewDuplicateKeyError is called to create an DuplicateKeyError error
func NewDuplicateKeyError(msg string) DuplicateKeyError {
	return DuplicateKeyError{msg: msg}
}

// DataCorruptionError indicates that the data is in an unexpected, illegal state
type DataCorruptionError struct {
	msg string
}

func (o DataCorruptionError) Error() string {
	return o.msg
}

// NewDataCorruptionError is called to create an DataCorruptionError error
func NewDataCorruptionError(msg string) DataCorruptionError {
	return DataCorruptionError{msg: msg}
}

// NotFoundError indicates that the entity that was expected to be affected by the operation was not found
type NotFoundError struct {
	msg string
	Err error
}

func (e NotFoundError) Error() string {
	return e.msg
}

// NewNotFoundError is called to create an NewNotFoundError error
func NewNotFoundError(msg string) NotFoundError {
	return NotFoundError{msg: msg}
}

const (
	// UniqueViolationError is raised if the unique constraint is violated
	UniqueViolationError = pq.ErrorCode("23505") // 'unique_violation'
)

// IsErrorCode a specific postgres specific error as defined by
// https://www.postgresql.org/docs/12/errcodes-appendix.html
func IsErrorCode(err error, errcode pq.ErrorCode) bool {
	var pgerr *pq.Error
	if errors.As(err, &pgerr) {
		return pgerr.Code == errcode
	}
	return false
}
