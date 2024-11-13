package ts

import (
	"bufio"
	"fmt"
	"slices"
	"strings"
	"time"

	"gabe565.com/moreutils/internal/cmdutil"
	"gabe565.com/moreutils/internal/util"
	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/must"
	"github.com/gravwell/gravwell/v3/timegrinder"
	"github.com/lestrrat-go/strftime"
	"github.com/spf13/cobra"
)

const (
	Name           = "ts"
	FlagMonotonic  = "monotonic"
	FlagIncrement  = "increment"
	FlagSinceStart = "since-start"
	FlagRelative   = "relative"
	FlagLocal      = "local"
	FlagMultiple   = "multiple"
)

func New(opts ...cobrax.Option) *cobra.Command {
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
	cmd.Flags().BoolP(FlagRelative, "r", false, "Convert existing timestamps from stdin to relative times")
	cmd.Flags().BoolP(FlagLocal, "l", false, "Parse to relative using local timezone instead of UTC")
	cmd.Flags().Bool(FlagMultiple, false, "Search lines for multiple timestamps when converting to relative. This is slower than the default behavior.")
	if err := cmd.Flags().MarkHidden(FlagMonotonic); err != nil {
		panic(err)
	}

	cmd.MarkFlagsMutuallyExclusive(FlagIncrement, FlagRelative)
	cmd.MarkFlagsMutuallyExclusive(FlagSinceStart, FlagRelative)

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func validArgs(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	var completions []string
	if must.Must2(cmd.Flags().GetBool(FlagIncrement)) || must.Must2(cmd.Flags().GetBool(FlagSinceStart)) {
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

	increment := must.Must2(cmd.Flags().GetBool(FlagIncrement))
	sinceStart := must.Must2(cmd.Flags().GetBool(FlagSinceStart))
	relative := must.Must2(cmd.Flags().GetBool(FlagRelative))
	parseLocal := must.Must2(cmd.Flags().GetBool(FlagLocal))
	multiple := must.Must2(cmd.Flags().GetBool(FlagMultiple))

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
	if relative {
		tg, err := timegrinder.New(timegrinder.Config{})
		if err != nil {
			return err
		}
		if parseLocal {
			tg.SetLocalTime()
		}
		for scanner.Scan() {
			line := scanner.Bytes()
			for offset := 0; offset < len(line); {
				ts, _, start, end, ok := tg.DebugMatch(line[offset:])
				if !ok {
					break
				}

				var replacement string
				if len(args) != 0 {
					replacement = formatter.FormatString(ts)
				} else if since := time.Since(ts).Round(time.Second); since < 0 {
					replacement = "in " + since.Abs().String()
				} else {
					replacement = since.String() + " ago"
				}

				line = slices.Concat(line[:start+offset], []byte(replacement), line[end+offset:])
				if !multiple {
					break
				}

				offset += len(replacement)
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", line); err != nil {
				return err
			}
		}
	} else {
		for scanner.Scan() {
			ts := time.Now()
			if increment || sinceStart {
				ts = time.Unix(0, 0).UTC().Add(time.Since(start))
				if increment {
					start = time.Now()
				}
			}

			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n",
				formatter.FormatString(ts),
				scanner.Bytes(),
			); err != nil {
				return err
			}
		}
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
