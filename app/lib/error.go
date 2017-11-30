package lib

import (
  "errors"
)

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}

func NewError(code int, message string) error {
  return StatusError{code, errors.New(message)}
}
