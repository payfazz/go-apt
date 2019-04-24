package httpError

import "net/http"

// TooManyRequestError is a struct to contain bad request http error
type TooManyRequestError struct {
	BaseError
}

// TooManyRequest is a constructor to create BadRequestError instance
func TooManyRequest(err error) HttpErrorInterface {
	return &TooManyRequestError{
		BaseError: code(http.StatusTooManyRequests, err),
	}
}
