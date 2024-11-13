package errno

import (
	"errors"

	"gabe565.com/moreutils/internal/cmdutil"
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra"
)

const (
	Name       = "errno"
	FlagList   = "list"
	FlagSearch = "search"
)

var (
	ErrUnknownNo = errors.New("unknown errno")
	ErrUnknown   = errors.New("unknown err name")
)

func New(opts ...cobrax.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " { name | code }",
		Short:   "Look up errno names and descriptions",
		RunE:    run,
		GroupID: cmdutil.Applet,
		Hidden:  !Supported,

		ValidArgsFunction: validArgs,
	}

	cmd.Flags().BoolP(FlagList, "l", false, "List all errno values")
	cmd.Flags().BoolP(FlagSearch, "s", false, "Search for errors whose description contains all the given words (case-insensitive)")

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}
