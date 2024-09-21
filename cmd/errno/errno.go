package errno

import (
	"github.com/gabe565/moreutils/internal/cmdutil"
	"github.com/spf13/cobra"
)

const (
	Name       = "errno"
	FlagList   = "list"
	FlagSearch = "search"
)

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " { name | code }",
		Short:   "Look up errno names and descriptions",
		Args:    cobra.MaximumNArgs(1),
		RunE:    run,
		GroupID: cmdutil.Applet,
		Hidden:  !Supported,
	}

	cmd.Flags().BoolP(FlagList, "l", false, "List all errno values")
	cmd.Flags().BoolP(FlagSearch, "s", false, "Search for errors whose description contains all the given words (case-insensitive)")

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}
