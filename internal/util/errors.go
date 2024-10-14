package util

import (
	"errors"
	"strconv"
)

var ErrNotAPipe = errors.New("this command should be run in a pipe")

func NewExitCodeError(code int) error {
	return &ExitCodeError{code}
}

type ExitCodeError struct {
	code int
}

func (e ExitCodeError) Error() string {
	return "exit code " + strconv.Itoa(e.code)
}

func (e ExitCodeError) ExitCode() int {
	return e.code
}
