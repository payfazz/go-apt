package httpError

import "net/http"

// UnauthorizedError is a struct to contain unauthorized http error
type UnauthorizedError struct {
	BaseError
}

// Unauthorized is a constructor to create UnauthorizedError instance
func Unauthorized(err error) *UnauthorizedError {
	return &UnauthorizedError{
		BaseError: code(http.StatusUnauthorized, err),
	}
}