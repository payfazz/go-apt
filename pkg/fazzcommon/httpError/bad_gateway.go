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

// IsBadGatewayError check whether given error is a BadGatewayError
func IsBadGatewayError(err error) bool {
	cause := getCause(err)
	_, ok := cause.(*BadGatewayError)
	return ok
}
