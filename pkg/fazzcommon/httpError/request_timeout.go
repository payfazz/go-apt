package httpError

import "net/http"

// RequestTimeoutError is a struct to contain request timeout http error
type RequestTimeoutError struct {
	BaseError
}

// RequestTimeout is a constructor to create NotFoundError instance
func RequestTimeout(err interface{}) error {
	return &RequestTimeoutError{
		BaseError: Code(http.StatusRequestTimeout, err),
	}
}

// IsRequestTimeoutError check whether given error is a RequestTimeoutError
func IsRequestTimeoutError(err error) bool {
	cause := getCause(err)
	_, ok := cause.(*RequestTimeoutError)
	return ok
}
