/**
 * @Author: lzw5399
 * @Date: 2021/3/27 13:33
 * @Desc:
 */
package util

import (
	"fmt"

	"github.com/pkg/errors"
)

// ErrorType is the type of an error
type ErrorType uint

const (
	NoType              ErrorType = iota
	BadRequest                    // 400
	NotFound                      // 404
	Forbidden                     // 403
	InternalServerError           // 500
)

type customError struct {
	errorType     ErrorType
	originalError error
	context       errorContext
}

type errorContext struct {
	Field   string
	Message string
}

// NewError creates a new customError
func (errorType ErrorType) New(err interface{}) error {
	var finalError error
	switch err.(type) {
	case error:
		finalError = err.(error)
	case string:
		finalError = errors.New(err.(string))
	}

	return customError{errorType: errorType, originalError: finalError}
}

// NewError creates a new customError with formatted message
func (errorType ErrorType) Newf(msg string, args ...interface{}) error {
	return customError{errorType: errorType, originalError: fmt.Errorf(msg, args...)}
}

// Error returns the message of a customError
func (error customError) Error() string {
	return error.originalError.Error()
}

// NewError creates a no type error
func NewError(err interface{}) error {
	var finalError error
	switch err.(type) {
	case error:
		finalError = err.(error)
	case string:
		finalError = errors.New(err.(string))
	}

	return customError{errorType: NoType, originalError: finalError}
}

// Newf creates a no type error with formatted message
func NewErrorf(msg string, args ...interface{}) error {
	return customError{errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.errorType
	}

	return NoType
}
