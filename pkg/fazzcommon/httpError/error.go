package httpError

import (
	"fmt"
	"net/http"

	"github.com/payfazz/go-errors"
)

// HttpErrorInterface is an interface for all http error
type HttpErrorInterface interface {
	GetCode() int
	GetStatusCode() string
	GetDetail() interface{}
}

// BaseError is a struct that contain basic requirements for http error struct
type BaseError struct {
	Code       int    `json:"-"`
	StatusCode string `json:"code"`
	Message    string `json:"message"`
}

// Error is a function to implement error interface
func (e *BaseError) Error() string {
	return fmt.Sprintf("%d %s: %s", e.Code, http.StatusText(e.Code), e.Message)
}

// GetCode is a function to return http error code
func (e *BaseError) GetCode() int {
	return e.Code
}

// GetStatusCode is a function to return error message
func (e *BaseError) GetStatusCode() string {
	return e.StatusCode
}

// GetDetail is a function to return raw error message
func (e *BaseError) GetDetail() interface{} {
	return e.Message
}

// Base is a constructor for http error with custom message
func Base(code int, message string, err interface{}) BaseError {
	return BaseError{
		Code:       code,
		StatusCode: message,
		Message:    getError(err),
	}
}

// Code is a constructor for http error with default status text message
func Code(code int, err interface{}) BaseError {
	return BaseError{
		Code:       code,
		StatusCode: http.StatusText(code),
		Message:    getError(err),
	}
}

// getError is a function for get error message from string or error
func getError(err interface{}) string {
	if v, ok := err.(error); ok {
		return v.Error()
	}
	return fmt.Sprint(err)
}

// getCause is a function to get error cause whether code use error interface or go-errors
func getCause(err error) error {
	if goErr, ok := err.(*errors.Error); ok {
		return goErr.Cause
	}

	return err
}
