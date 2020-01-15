package httpError

import "net/http"

// BadRequestError is a struct to contain bad request http error
type BadRequestError struct {
	BaseError
}

// BadRequest is a constructor to create BadRequestError instance
func BadRequest(err interface{}) error {
	return &BadRequestError{
		BaseError: Code(http.StatusBadRequest, err),
	}
}

// IsBadRequestError check whether given error is a BadRequestError
func IsBadRequestError(err error) bool {
	cause := getCause(err)
	_, ok := cause.(*BadRequestError)
	return ok
}
