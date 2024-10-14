package main

import (
	"errors"
	"os"
	"os/exec"

	"gabe565.com/moreutils/cmd"
	"gabe565.com/moreutils/internal/cmdutil"
	"gabe565.com/moreutils/internal/util"
)

var version = "beta"

func main() {
	root := cmd.New(os.Args[0], cmdutil.WithVersion(version))
	root.SilenceErrors = true
	if err := root.Execute(); err != nil {
		var execErr *exec.ExitError
		if errors.As(err, &execErr) {
			os.Exit(execErr.ExitCode())
		}

		var exitCodeErr *util.ExitCodeError
		if errors.As(err, &exitCodeErr) {
			os.Exit(exitCodeErr.ExitCode())
		}

		root.PrintErrln(root.ErrPrefix(), err.Error())
		os.Exit(1)
	}
}
