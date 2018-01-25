package lib

import (
  "errors"
  "github.com/eriklindqvist/recepies_auth/log"
  "fmt"
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
  log.Err(fmt.Sprintf("%d %s", code, message))
  return StatusError{code, errors.New(message)}
}
