package httpError

import "net/http"

type BadRequestError struct {
	BaseError
}

func BadRequest(err error) *BadRequestError {
	return &BadRequestError{
		BaseError: code(http.StatusBadRequest, err),
	}
}