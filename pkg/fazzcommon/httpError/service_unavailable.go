package httpError

import "net/http"

// ServiceUnavailableError is a struct to contain service unavailable http error
type ServiceUnavailableError struct {
	BaseError
}

// ServiceUnavailable is a constructor to create NotFoundError instance
func ServiceUnavailable(err interface{}) error {
	return &ServiceUnavailableError{
		BaseError: Code(http.StatusServiceUnavailable, err),
	}
}
