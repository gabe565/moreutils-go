//go:build !unix

package errno

import (
	"errors"
	"runtime"

	"github.com/spf13/cobra"
)

const Supported = false

var ErrUnsupported = errors.New(Name + " is unsupported on " + runtime.GOOS)

func run(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true
	return ErrUnsupported
}
