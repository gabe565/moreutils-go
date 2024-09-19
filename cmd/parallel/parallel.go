package parallel

import (
	"errors"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/gabe565/moreutils/cmd/cmdutil"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

const (
	Name        = "parallel"
	FlagJobs    = "jobs"
	FlagReplace = "replace"
	FlagNumArgs = "num-args"
)

func New(opts ...cmdutil.Option) *cobra.Command {
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
	cmd.Flags().BoolP(FlagReplace, "i", false, `Normally the argument is added to the end of the command. With this option, instances of "{}" in the command are replaced with the argument.`)
	cmd.Flags().IntP(FlagNumArgs, "n", 1, "Number of arguments to pass to a command at a time. Default is 1. Incompatible with -i")
	cmd.MarkFlagsMutuallyExclusive(FlagReplace, FlagNumArgs)

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

	numArgs, err := cmd.Flags().GetInt(FlagNumArgs)
	if err != nil {
		panic(err)
	}

	replace, err := cmd.Flags().GetBool(FlagReplace)
	if err != nil {
		panic(err)
	}

	var group errgroup.Group
	group.SetLimit(numJobs)

	execCmd := args[:sepIdx]
	for args := range slices.Chunk(args[sepIdx+1:], numArgs) {
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
	numJobsStr, err := cmd.Flags().GetString(FlagJobs)
	if err != nil {
		panic(err)
	}

	var jobs int
	if strings.HasSuffix(numJobsStr, "%") {
		pct, err := strconv.Atoi(strings.TrimSuffix(numJobsStr, "%"))
		if err != nil {
			return 0, err
		}

		jobs = runtime.NumCPU() * pct / 100
	} else {
		jobs, err = strconv.Atoi(numJobsStr)
		if err != nil {
			return 0, err
		}
	}

	return jobs, nil
}
