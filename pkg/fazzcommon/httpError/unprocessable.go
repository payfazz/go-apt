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

// IsUnprocessableEntityError check whether given error is a UnprocessableEntityError
func IsUnprocessableEntityError(err error) bool {
	cause := getCause(err)
	_, ok := cause.(*UnprocessableEntityError)
	return ok
}
