package httpError

import "net/http"

type ForbiddenError struct {
	BaseError
}

func Forbidden(err error) *ForbiddenError {
	return &ForbiddenError{
		BaseError: code(http.StatusForbidden, err),
	}
}
