package ts

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"github.com/gabe565/moreutils/internal/cmdutil"
	"github.com/gabe565/moreutils/internal/util"
	"github.com/lestrrat-go/strftime"
	"github.com/spf13/cobra"
)

const (
	Name           = "ts"
	FlagMonotonic  = "monotonic"
	FlagIncrement  = "increment"
	FlagSinceStart = "since-start"
)

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " [format]",
		Short:   "Timestamp standard input",
		Args:    cobra.MaximumNArgs(1),
		RunE:    run,
		GroupID: cmdutil.Applet,

		ValidArgsFunction: validArgs,
	}

	cmd.Flags().BoolP(FlagMonotonic, "m", false, "Use the system's monotonic clock")
	cmd.Flags().BoolP(FlagIncrement, "i", false, "Timestamps will be the time elapsed since the last log")
	cmd.Flags().BoolP(FlagSinceStart, "s", false, "Timestamps will be the time elapsed since start of the program")
	if err := cmd.Flags().MarkHidden(FlagMonotonic); err != nil {
		panic(err)
	}

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func validArgs(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	var completions []string
	if util.Must2(cmd.Flags().GetBool(FlagIncrement)) || util.Must2(cmd.Flags().GetBool(FlagSinceStart)) {
		completions = append(completions,
			"%T",
			"%.T",
		)
	} else {
		completions = []string{
			"%b %e %T",
			"%b %e %.T",
			"%a %b %e %T %Y",
			"%F %T",
			"%F %.T",
			"%FT%T%z",
		}
	}
	now := time.Now()
	for i, completion := range completions {
		formatter, err := newFormatter(completion)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		completions[i] += "\t" + formatter.FormatString(now)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveKeepOrder
}

func run(cmd *cobra.Command, args []string) error {
	if util.IsTerminal(cmd.InOrStdin()) {
		return util.ErrNotAPipe
	}
	cmd.SilenceUsage = true

	increment := util.Must2(cmd.Flags().GetBool(FlagIncrement))
	sinceStart := util.Must2(cmd.Flags().GetBool(FlagSinceStart))

	format := "%Y-%m-%d %H:%M:%S"
	switch {
	case len(args) != 0:
		format = args[0]
	case increment, sinceStart:
		format = "%H:%M:%S"
	}
	formatter, err := newFormatter(format)
	if err != nil {
		return err
	}

	start := time.Now()
	scanner := bufio.NewScanner(cmd.InOrStdin())
	for scanner.Scan() {
		ts := time.Now()
		if increment || sinceStart {
			ts = time.Unix(0, 0).UTC().Add(time.Since(start))
			if increment {
				start = time.Now()
			}
		}
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n",
			formatter.FormatString(ts),
			scanner.Bytes(),
		)
	}
	return scanner.Err()
}

func newFormatter(src string) (*strftime.Strftime, error) {
	src = strings.ReplaceAll(src, "%.S", "%S.%f")
	src = strings.ReplaceAll(src, "%.T", "%T.%f")
	src = strings.ReplaceAll(src, "%.s", "%s.%f")

	return strftime.New(src,
		strftime.WithMilliseconds('L'),
		strftime.WithMicroseconds('f'),
		strftime.WithUnixSeconds('s'),
	)
}
