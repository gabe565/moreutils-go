package ts

import (
	"bufio"
	"fmt"
	"time"

	"github.com/gabe565/moreutils/cmd/cmdutil"
	"github.com/gabe565/moreutils/internal/util"
	"github.com/ncruces/go-strftime"
	"github.com/spf13/cobra"
)

const (
	Name          = "ts"
	FlagMonotonic = "monotonic"
)

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " [format]",
		Short:   "Timestamp standard input",
		Args:    cobra.NoArgs,
		RunE:    run,
		GroupID: cmdutil.Applet,
	}

	cmd.Flags().BoolP(FlagMonotonic, "m", false, "Use the system's monotonic clock")
	if err := cmd.Flags().MarkDeprecated(FlagMonotonic, "the monotonic clock is always used"); err != nil {
		panic(err)
	}

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if util.IsTerminal(cmd.InOrStdin()) {
		return cmd.Usage()
	}
	cmd.SilenceUsage = true

	format := time.DateTime
	if len(args) > 0 {
		var err error
		if format, err = strftime.Layout(args[0]); err != nil {
			return err
		}
	}

	scanner := bufio.NewScanner(cmd.InOrStdin())
	for scanner.Scan() {
		fmt.Println(timestamp(format), scanner.Text())
	}
	return scanner.Err()
}

func timestamp(format string) string {
	return time.Now().Format(format)
}
