package httpError

import "net/http"

// MethodNotAllowedError is a struct to contain method not allowed http error
type MethodNotAllowedError struct {
	BaseError
}

// MethodNotAllowed is a constructor to create NotFoundError instance
func MethodNotAllowed(err interface{}) error {
	return &MethodNotAllowedError{
		BaseError: Code(http.StatusMethodNotAllowed, err),
	}
}

// IsMethodNotAllowedError check whether given error is a MethodNotAllowedError
func IsMethodNotAllowedError(err error) bool {
	cause := getCause(err)
	_, ok := cause.(*MethodNotAllowedError)
	return ok
}
