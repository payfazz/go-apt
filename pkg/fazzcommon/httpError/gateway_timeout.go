package httpError

import "net/http"

// GatewayTimeoutError is a struct to contain gateway timeout http error
type GatewayTimeoutError struct {
	BaseError
}

// GatewayTimeout is a constructor to create NotFoundError instance
func GatewayTimeout(err interface{}) error {
	return &GatewayTimeoutError{
		BaseError: Code(http.StatusGatewayTimeout, err),
	}
}
