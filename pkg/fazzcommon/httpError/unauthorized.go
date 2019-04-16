package httpError

import "net/http"

type UnauthorizedError struct {
	BaseError
}

func Unauthorized(err error) *UnauthorizedError {
	return &UnauthorizedError{
		BaseError: code(http.StatusUnauthorized, err),
	}
}