package ecode

import (
	"errors"
	"fmt"
)

type Error struct {
	code  Code
	data  any
	cause error
}

func New(code Code) *Error { return &Error{code: code} }

func Wrap(code Code, cause error) *Error { return &Error{code: code, cause: cause} }

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.ErrorMessage(), e.cause)
	}
	return e.ErrorMessage()
}

func (e *Error) Unwrap() error { return e.cause }

func (e *Error) Code() Code { return e.code }

func (e *Error) HTTPStatus() int { return HTTPStatus(e.code) }

func (e *Error) ErrorMessage() string {
	return Message(e.code)
}

func (e *Error) Data() any { return e.data }

func (e *Error) WithData(data any) *Error {
	e.data = data
	return e
}

func FromError(err error) (*Error, bool) {
	var target *Error
	ok := errors.As(err, &target)
	return target, ok
}
