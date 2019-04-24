package httpError

import "net/http"

// InternalServerError is a struct to contain internal server http error
type InternalServerError struct {
	BaseError
}

// InternalServer is a constructor to create InternalServerError instance
func InternalServer(err error) HttpErrorInterface {
	return &InternalServerError{
		BaseError: code(http.StatusInternalServerError, err),
	}
}
