package pee

import (
	"io"
	"os/exec"

	"github.com/gabe565/moreutils/internal/cmdutil"
	"github.com/gabe565/moreutils/internal/util"
	"github.com/spf13/cobra"
)

const (
	Name                  = "pee"
	FlagIgnoreSigpipe     = "ignore-sigpipe"
	FlagIgnoreWriteErrors = "ignore-write-errors"
)

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " command...",
		Short:   "Tee standard input to pipes",
		Args:    cobra.MinimumNArgs(1),
		RunE:    run,
		GroupID: cmdutil.Applet,
	}

	cmd.Flags().Bool(FlagIgnoreSigpipe, true, "Ignores sigpipe")
	cmd.Flags().Bool(FlagIgnoreWriteErrors, true, "Ignores write errors")
	if err := cmd.Flags().MarkHidden(FlagIgnoreSigpipe); err != nil {
		panic(err)
	}
	if err := cmd.Flags().MarkHidden(FlagIgnoreWriteErrors); err != nil {
		panic(err)
	}

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	cmds := make([]*exec.Cmd, 0, len(args))
	pipes := make([]io.Writer, 0, len(args))
	pipeClosers := make([]io.WriteCloser, 0, len(args))
	var errs []error
	for _, arg := range args {
		e := exec.Command("sh", "-c", arg)
		e.Stdout = cmd.OutOrStdout()
		e.Stderr = cmd.ErrOrStderr()
		stdin, err := e.StdinPipe()
		if err != nil {
			return err
		}

		if err := e.Start(); err != nil {
			errs = append(errs, err)
		}

		cmds = append(cmds, e)
		pipes = append(pipes, util.NewSuppressErrorWriter(stdin))
		pipeClosers = append(pipeClosers, stdin)
	}

	if _, err := io.Copy(io.MultiWriter(pipes...), cmd.InOrStdin()); err != nil {
		return err
	}

	for i, e := range cmds {
		if err := pipeClosers[i].Close(); err != nil {
			errs = append(errs, err)
		}
		if err := e.Wait(); err != nil {
			errs = append(errs, err)
		}
	}

	return util.JoinErrors(errs...)
}
