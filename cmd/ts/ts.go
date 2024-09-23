package ts

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"github.com/gabe565/moreutils/internal/cmdutil"
	"github.com/gabe565/moreutils/internal/util"
	"github.com/ncruces/go-strftime"
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
			"%H:%M:%S",
			"%H:%M:%.S",
		)
	} else {
		completions = []string{
			"%b %e %H:%M:%S",
			"%b %e %H:%M:%.S",
			"%a %b %e %H:%M:%S %Y",
			"%Y-%m-%d %H:%M:%S",
			"%Y-%m-%d %H:%M:%.S",
			"%Y-%m-%dT%H:%M:%S%z",
		}
	}
	now := time.Now()
	for i, completion := range completions {
		format, err := convertFormat(completion)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		completions[i] += "\t" + now.Format(format)
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

	format := time.DateTime
	switch {
	case len(args) != 0:
		var err error
		format, err = convertFormat(args[0])
		if err != nil {
			return err
		}
	case increment, sinceStart:
		format = time.TimeOnly
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
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", ts.Format(format), scanner.Bytes())
	}
	return scanner.Err()
}

func convertFormat(src string) (string, error) {
	const secWithMs = "\x1bSEC_WITH_MS"

	src = strings.ReplaceAll(src, "%.S", secWithMs)
	src = strings.ReplaceAll(src, "%.T", "%H:%M:"+secWithMs)

	var err error
	if src, err = strftime.Layout(src); err != nil {
		return "", err
	}

	src = strings.ReplaceAll(src, secWithMs, "05.000000")
	return src, nil
}
