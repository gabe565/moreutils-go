//go:build !unix

package errno

import (
	"errors"
	"runtime"

	"github.com/spf13/cobra"
)

const Supported = false

var ErrUnsupported = errors.New(Name + " is unsupported on " + runtime.GOOS)

func validArgs(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveDefault
}

func run(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true
	return ErrUnsupported
}
