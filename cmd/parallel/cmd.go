package parallel

import (
	"context"
	"errors"
	"log/slog"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"gabe565.com/moreutils/internal/cmdutil"
	"gabe565.com/moreutils/internal/loadavg"
	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/must"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

const (
	Name        = "parallel"
	FlagJobs    = "jobs"
	FlagLoad    = "load"
	FlagReplace = "replace"
	FlagNumArgs = "num-args"
)

func New(opts ...cobrax.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " [flags] command -- arg...",
		Short:   "Run multiple jobs at once",
		Args:    cobra.MinimumNArgs(1),
		RunE:    run,
		GroupID: cmdutil.Applet,

		DisableFlagsInUseLine: true,
	}

	cmd.Flags().SetInterspersed(false)
	cmd.Flags().StringP(FlagJobs, "j", strconv.Itoa(runtime.NumCPU()), "Number of jobs to run in parallel. Can be a number or a percentage of CPU cores.")
	cmd.Flags().Float64P(FlagLoad, "l", 0, "Wait until the system's load average is below a limit before starting jobs")
	cmd.Flags().BoolP(FlagReplace, "i", false, `Normally the argument is added to the end of the command. With this option, instances of "{}" in the command are replaced with the argument.`)
	cmd.Flags().IntP(FlagNumArgs, "n", 1, "Number of arguments to pass to a command at a time. Default is 1. Incompatible with -i")
	cmd.MarkFlagsMutuallyExclusive(FlagReplace, FlagNumArgs)

	if !loadavg.Supported {
		if err := cmd.Flags().MarkHidden(FlagLoad); err != nil {
			panic(err)
		}
	}

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

var ErrMissingSeparator = errors.New("missing separator")

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	sepIdx := slices.Index(args, "--")
	if sepIdx == -1 {
		return ErrMissingSeparator
	}

	numJobs, err := parseNumJobs(cmd)
	if err != nil {
		return err
	}

	numArgs := must.Must2(cmd.Flags().GetInt(FlagNumArgs))
	replace := must.Must2(cmd.Flags().GetBool(FlagReplace))

	maxLoad := must.Must2(cmd.Flags().GetFloat64(FlagLoad))
	loadAvg := loadavg.New()

	var group errgroup.Group
	group.SetLimit(numJobs)

	execCmd := args[:sepIdx]
	for args := range slices.Chunk(args[sepIdx+1:], numArgs) {
		if maxLoad != 0 {
			if err := loadAvg.WaitBelow(context.Background(), maxLoad, time.Second); err != nil {
				slog.Error("Failed to determine loadavg. Try again without -l")
				return err
			}
		}

		group.Go(func() error {
			execCmd := buildCmd(execCmd, args, replace)
			e := exec.Command(execCmd[0], execCmd[1:]...)
			e.Stdin = cmd.InOrStdin()
			e.Stdout = cmd.OutOrStdout()
			e.Stderr = cmd.OutOrStderr()
			return e.Run()
		})
	}
	return group.Wait()
}

func buildCmd(args []string, arg []string, replace bool) []string {
	args = slices.Clone(args)
	if replace {
		for i, v := range args {
			args[i] = strings.ReplaceAll(v, "{}", arg[0])
		}
	} else {
		args = append(args, arg...)
	}
	return args
}

func parseNumJobs(cmd *cobra.Command) (int, error) {
	numJobsStr := must.Must2(cmd.Flags().GetString(FlagJobs))
	var jobs int
	if strings.HasSuffix(numJobsStr, "%") {
		pct, err := strconv.Atoi(strings.TrimSuffix(numJobsStr, "%"))
		if err != nil {
			return 0, err
		}

		jobs = runtime.NumCPU() * pct / 100
	} else {
		var err error
		if jobs, err = strconv.Atoi(numJobsStr); err != nil {
			return 0, err
		}
	}

	return jobs, nil
}
