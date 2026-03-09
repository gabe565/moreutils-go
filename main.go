package main

import (
	"errors"
	"os"
	"os/exec"

	"gabe565.com/moreutils/cmd"
	"gabe565.com/moreutils/internal/util"
	"gabe565.com/utils/cobrax"
)

var version = "beta"

func main() {
	root := cmd.New(os.Args[0], cobrax.WithVersion(version))
	root.SilenceErrors = true
	if err := root.Execute(); err != nil {
		if err, ok := errors.AsType[*exec.ExitError](err); ok {
			os.Exit(err.ExitCode())
		}

		if err, ok := errors.AsType[*util.ExitCodeError](err); ok {
			os.Exit(err.ExitCode())
		}

		root.PrintErrln(root.ErrPrefix(), err.Error())
		os.Exit(1)
	}
}
