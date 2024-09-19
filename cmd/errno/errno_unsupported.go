//go:build !(aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos)

package errno

import (
	"errors"
	"runtime"

	"github.com/gabe565/moreutils/internal/cmdutil"
	"github.com/spf13/cobra"
)

const (
	Name      = "errno"
	Supported = false
)

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name,
		Short:   "Look up errno names and descriptions. Not supported on " + runtime.GOOS + ".",
		RunE:    run,
		GroupID: cmdutil.Applet,
	}
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

var ErrUnsupported = errors.New("this command is not supported on " + runtime.GOOS)

func run(_ *cobra.Command, _ []string) error {
	return ErrUnsupported
}
