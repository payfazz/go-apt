package httpError

import "net/http"

// TooManyRequestError is a struct to contain bad request http error
type TooManyRequestError struct {
	BaseError
}

// TooManyRequest is a constructor to create BadRequestError instance
func TooManyRequest(err interface{}) error {
	return &TooManyRequestError{
		BaseError: Code(http.StatusTooManyRequests, err),
	}
}

// IsTooManyRequestError check whether given error is a TooManyRequestError
func IsTooManyRequestError(err error) bool {
	cause := getCause(err)
	_, ok := cause.(*TooManyRequestError)
	return ok
}
