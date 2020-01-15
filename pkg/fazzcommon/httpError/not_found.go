package httpError

import "net/http"

// NotFoundError is a struct to contain not found http error
type NotFoundError struct {
	BaseError
}

// NotFound is a constructor to create NotFoundError instance
func NotFound(err interface{}) error {
	return &NotFoundError{
		BaseError: Code(http.StatusNotFound, err),
	}
}

// IsNotFoundError check whether given error is a NotFoundError
func IsNotFoundError(err error) bool {
	cause := getCause(err)
	_, ok := cause.(*NotFoundError)
	return ok
}
