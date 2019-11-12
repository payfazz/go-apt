package httpError

import "net/http"

// NotImplementedError is a struct to contain not implemented http error
type NotImplementedError struct {
	BaseError
}

// NotImplemented is a constructor to create NotFoundError instance
func NotImplemented(err interface{}) error {
	return &NotImplementedError{
		BaseError: Code(http.StatusNotImplemented, err),
	}
}
