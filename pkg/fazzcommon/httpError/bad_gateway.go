package httpError

import "net/http"

// BadGatewayError is a struct to contain bad gateway http error
type BadGatewayError struct {
	BaseError
}

// BadGateway is a constructor to create NotFoundError instance
func BadGateway(err interface{}) error {
	return &BadGatewayError{
		BaseError: Code(http.StatusBadGateway, err),
	}
}
