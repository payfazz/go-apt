package httpError

import (
	"fmt"
	"net/http"
)

// HttpErrorInterface is an interface for all http error
type HttpErrorInterface interface {
	GetCode() int
	GetMessage() string
	GetDetail() interface{}
}

// BaseError is a struct that contain basic requirements for http error struct
type BaseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Trace   string `json:"error"`
}

// Error is a function to implement error interface
func (e *BaseError) Error() string {
	return fmt.Sprintf("%d %s: %s", e.Code, http.StatusText(e.Code), e.Trace)
}

// GetCode is a function to return http error code
func (e *BaseError) GetCode() int {
	return e.Code
}

// GetMessage is a function to return error message
func (e *BaseError) GetMessage() string {
	return e.Message
}

// GetDetail is a function to return raw error message
func (e *BaseError) GetDetail() interface{} {
	return e.Trace
}

// base is a constructor for http error with custom message
func base(code int, message string, err error) BaseError {
	return BaseError{
		Code:    code,
		Message: message,
		Trace:   err.Error(),
	}
}

// code is a constructor for http error with default status text message
func code(code int, err string) BaseError {
	return BaseError{
		Code:    code,
		Message: http.StatusText(code),
		Trace:   err,
	}
}
