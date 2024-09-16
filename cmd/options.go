package cmd

import (
	"github.com/gabe565/moreutils/cmd/cmdutil"
	"github.com/spf13/cobra"
)

func WithVersion(version string) cmdutil.Option {
	return func(cmd *cobra.Command) {
		cmd.Version = buildVersion(version)
		cmd.InitDefaultVersionFlag()
	}
}
