package httpError

import "net/http"

// ForbiddenError is a struct to contain forbidden http error
type ForbiddenError struct {
	BaseError
}

// Forbidden is a constructor to create ForbiddenError instance
func Forbidden(err interface{}) error {
	return &ForbiddenError{
		BaseError: Code(http.StatusForbidden, err),
	}
}

// IsForbiddenError check whether given error is a ForbiddenError
func IsForbiddenError(err error) bool {
	cause := getCause(err)
	_, ok := cause.(*ForbiddenError)
	return ok
}
