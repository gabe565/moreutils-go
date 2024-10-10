//go:build unix

package errno

import (
	"strings"
	"syscall"

	"golang.org/x/exp/constraints"
	"golang.org/x/sys/unix"
)

type Option func(e *Errno)

func WithName(name string) Option {
	return func(e *Errno) {
		e.name = name
	}
}

func WithDescription(description string) Option {
	return func(e *Errno) {
		e.desc = description
	}
}

func New[T unix.Errno | constraints.Integer](unixErr T, options ...Option) *Errno {
	err := &Errno{Errno: unix.Errno(unixErr)}
	for _, opt := range options {
		opt(err)
	}
	return err
}

type Errno struct { //nolint:errname
	syscall.Errno
	name, desc string
}

func (e *Errno) Valid() bool {
	return e.Name() != ""
}

func (e *Errno) Name() string {
	if e.name == "" {
		e.name = unix.ErrnoName(e.Errno)
	}
	return e.name
}

func (e *Errno) Error() string {
	if e.desc == "" {
		e.desc = e.Errno.Error()
		if len(e.desc) >= 2 {
			e.desc = strings.ToUpper(e.desc[0:1]) + e.desc[1:]
		}
	}
	return e.desc
}
