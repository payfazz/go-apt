package httpError

import "net/http"

// GoneError is a struct to contain bad request http error
type GoneError struct {
	BaseError
}

// Gone is a constructor to create GoneError instance
func Gone(err interface{}) error {
	return &GoneError{
		BaseError: Code(http.StatusGone, err),
	}
}
