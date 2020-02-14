package v1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CheckErrorCode returns true, if the given err is not nil,
// is of type grpc Status and the code equals the given code.
func CheckErrorCode(err error, code codes.Code) bool {
	if err == nil {
		return false
	}
	st, ok := status.FromError(err)
	if !ok {
		return false
	}
	return st.Code() == code
}

// IsNotFound checks if the given error is a notfound error.
func IsNotFound(err error) bool {
	return CheckErrorCode(err, codes.NotFound)
}

// IsConflict checks if the given error is a conflict error.
// Example: key already exists on create.
func IsConflict(err error) bool {
	return CheckErrorCode(err, codes.AlreadyExists)
}

// IsInternal checks if the given error is an internal server error.
func IsInternal(err error) bool {
	return CheckErrorCode(err, codes.Internal)
}

// IsOptimistickLockError checks if the given error is a Optimistic Lock Error,
// which indicates that you read an entity and tried to update this entity,
// but it changed in the datastore by another party in the meantime.
func IsOptimistickLockError(err error) bool {
	return CheckErrorCode(err, codes.FailedPrecondition)
}
