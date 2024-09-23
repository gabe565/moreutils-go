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
	Name          = "ts"
	FlagMonotonic = "monotonic"
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
	if err := cmd.Flags().MarkHidden(FlagMonotonic); err != nil {
		panic(err)
	}

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func validArgs(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	completions := []string{
		"%b %e %H:%M:%S",
		"%b %e %H:%M:%.S",
		"%a %b %e %H:%M:%S %Y",
		"%Y-%m-%d %H:%M:%S",
		"%Y-%m-%d %H:%M:%.S",
		"%Y-%m-%dT%H:%M:%S%z",
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

	format := time.DateTime
	if len(args) != 0 {
		var err error
		format, err = convertFormat(args[0])
		if err != nil {
			return err
		}
	}

	scanner := bufio.NewScanner(cmd.InOrStdin())
	for scanner.Scan() {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", time.Now().Format(format), scanner.Bytes())
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
