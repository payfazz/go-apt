package httpError

import "net/http"

// UnauthorizedError is a struct to contain unauthorized http error
type UnauthorizedError struct {
	BaseError
}

// Unauthorized is a constructor to create UnauthorizedError instance
func Unauthorized(err interface{}) error {
	return &UnauthorizedError{
		BaseError: Code(http.StatusUnauthorized, err),
	}
}

// IsUnauthorizedError check whether given error is a UnauthorizedError
func IsUnauthorizedError(err error) bool {
	cause := getCause(err)
	_, ok := cause.(*UnauthorizedError)
	return ok
}
