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
