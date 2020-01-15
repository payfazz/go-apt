package httpError

import "net/http"

// InternalServerError is a struct to contain internal server http error
type InternalServerError struct {
	BaseError
}

// InternalServer is a constructor to create InternalServerError instance
func InternalServer(err interface{}) error {
	return &InternalServerError{
		BaseError: Code(http.StatusInternalServerError, err),
	}
}

// IsInternalServerError check whether given error is a InternalServerError
func IsInternalServerError(err error) bool {
	cause := getCause(err)
	_, ok := cause.(*InternalServerError)
	return ok
}
