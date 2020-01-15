package httpError

import "net/http"

// InsufficientStorageError is a struct to contain insufficient storage http error
type InsufficientStorageError struct {
	BaseError
}

// InsufficientStorage is a constructor to create NotFoundError instance
func InsufficientStorage(err interface{}) error {
	return &InsufficientStorageError{
		BaseError: Code(http.StatusInsufficientStorage, err),
	}
}

// IsInsufficientStorageError check whether given error is a InsufficientStorageError
func IsInsufficientStorageError(err error) bool {
	cause := getCause(err)
	_, ok := cause.(*InsufficientStorageError)
	return ok
}
