package services

import (
	"errors"
	"fmt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type BaseValidationError struct {
	msg string
}

func (e *BaseValidationError) Error() string {
	return e.msg
}

func NewValidationError(msg string) *BaseValidationError {
	return &BaseValidationError{msg: msg}
}

type SpecificValidationError struct {
	*BaseValidationError
	field string
}

func NewSpecificValidationError(field, msg string) *SpecificValidationError {
	return &SpecificValidationError{
		BaseValidationError: NewValidationError(msg),
		field:               field,
	}
}

func (e *SpecificValidationError) Error() string {
	return fmt.Sprintf("field '%s': %s", e.field, e.msg)
}

func (e *SpecificValidationError) Unwrap() error {
	return e.BaseValidationError
}
