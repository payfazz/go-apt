package httpError

import "net/http"

type InternalServerError struct {
	BaseError
}

func InternalServer(err error) *InternalServerError {
	return &InternalServerError{
		BaseError: code(http.StatusInternalServerError, err),
	}
}
