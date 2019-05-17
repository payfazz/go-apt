package httpError

import "net/http"

// UnprocessableEntityError is a struct to contain bad request http error
type UnprocessableEntityError struct {
	BaseError
}

// UnprocessableEntity is a constructor to create UnprocessableEntityError instance
func UnprocessableEntity(err interface{}) error {
	return &UnprocessableEntityError{
		BaseError: Code(http.StatusUnprocessableEntity, err),
	}
}
