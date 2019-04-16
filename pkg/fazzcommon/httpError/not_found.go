package httpError

import "net/http"

type NotFoundError struct {
	BaseError
}

func NotFound(err error) *NotFoundError {
	return &NotFoundError{
		BaseError: code(http.StatusNotFound, err),
	}
}
