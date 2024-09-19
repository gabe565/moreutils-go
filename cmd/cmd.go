package cmd

import (
	"github.com/gabe565/moreutils/cmd/install"
	"github.com/gabe565/moreutils/internal/cmdutil"
	"github.com/gabe565/moreutils/internal/cmdutil/subcommands"
	"github.com/spf13/cobra"
)

const Name = "moreutils"

func New(name string, opts ...cmdutil.Option) *cobra.Command {
	if name != Name {
		if cmd, err := subcommands.Choose(name, opts...); err == nil {
			return cmd
		}
	}

	cmd := &cobra.Command{
		Use:   Name,
		Short: "A collection of the Unix tools that nobody thought to write long ago when Unix was young",

		DisableAutoGenTag: true,
	}
	cmd.AddGroup(&cobra.Group{
		ID:    cmdutil.Applet,
		Title: "Applets:",
	})
	cmd.AddCommand(install.New())
	cmd.AddCommand(subcommands.All()...)
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}
