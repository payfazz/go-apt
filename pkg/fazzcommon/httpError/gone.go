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

// IsGoneError check whether given error is a GoneError
func IsGoneError(err error) bool {
	cause := getCause(err)
	_, ok := cause.(*GoneError)
	return ok
}
