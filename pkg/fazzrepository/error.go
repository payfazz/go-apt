package fazzrepository

import "github.com/payfazz/go-errors"

// emptyResultError appear when data returned as empty
type emptyResultError struct{}

// Error return error text
func (e *emptyResultError) Error() string {
	return "no rows in result set"
}

// NewEmptyResultError creates new instance of EmptyResultError
func NewEmptyResultError() error {
	return &emptyResultError{}
}

// IsEmptyResult check if instance of error is EmptyResultError
func IsEmptyResult(it interface{}) bool {
	switch it.(type) {
	case *errors.Error:
		return IsEmptyResult(it.(*errors.Error).Cause)
	case emptyResultError:
		return true
	case *emptyResultError:
		return true
	default:
		return false
	}
}
