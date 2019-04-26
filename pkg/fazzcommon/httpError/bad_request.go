package httpError

import "net/http"

// BadRequestError is a struct to contain bad request http error
type BadRequestError struct {
	BaseError
}

// BadRequest is a constructor to create BadRequestError instance
func BadRequest(err interface{}) error {
	return &BadRequestError{
		BaseError: code(http.StatusBadRequest, err),
	}
}