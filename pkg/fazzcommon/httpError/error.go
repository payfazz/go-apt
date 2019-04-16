package httpError

import (
	"fmt"
	"net/http"
)

type HttpErrorInterface interface {
	GetCode() int
	GetMessage() string
	GetDetail() interface{}
}

type BaseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Trace   string `json:"error"`
}

func (e *BaseError) Error() string {
	return fmt.Sprintf("%d %s: %s", e.Code, http.StatusText(e.Code), e.Trace)
}

func (e *BaseError) GetCode() int {
	return e.Code
}

func (e *BaseError) GetMessage() string {
	return e.Message
}

func (e *BaseError) GetDetail() interface{} {
	return e.Trace
}

func base(code int, message string, err error) BaseError {
	return BaseError{
		Code:    code,
		Message: message,
		Trace:   err.Error(),
	}
}

func code(code int, err error) BaseError {
	return BaseError{
		Code:    code,
		Message: http.StatusText(code),
		Trace:   err.Error(),
	}
}
