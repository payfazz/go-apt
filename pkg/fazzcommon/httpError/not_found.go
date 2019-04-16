package httpError

import "net/http"

// NotFoundError is a struct to contain not found http error
type NotFoundError struct {
	BaseError
}

// NotFound is a constructor to create NotFoundError instance
func NotFound(err error) *NotFoundError {
	return &NotFoundError{
		BaseError: code(http.StatusNotFound, err),
	}
}
