package pee

import (
	"io"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gabe565/moreutils/internal/cmdutil"
	"github.com/gabe565/moreutils/internal/util"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const (
	Name = "pee"

	FlagIgnoreSigpipe       = "ignore-sigpipe"
	FlagNoIgnoreSigpipe     = "no-ignore-sigpipe"
	FlagIgnoreWriteErrors   = "ignore-write-errors"
	FlagNoIgnoreWriteErrors = "no-ignore-write-errors"
)

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " command...",
		Short:   "Tee standard input to pipes",
		Args:    cobra.MinimumNArgs(1),
		RunE:    run,
		GroupID: cmdutil.Applet,
	}

	cmd.Flags().Bool(FlagIgnoreSigpipe, true, "")
	cmd.Flags().Bool(FlagIgnoreWriteErrors, true, "")
	cmd.Flags().Bool(FlagNoIgnoreSigpipe, false, "Do not ignore write errors")
	cmd.Flags().Bool(FlagNoIgnoreWriteErrors, false, "Do not ignore SIGPIPE errors")
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

	ignoreSigpipe, ignoreWriteErrs := true, true
	cmd.Flags().Visit(func(f *flag.Flag) {
		switch f.Name {
		case FlagIgnoreSigpipe:
			ignoreSigpipe = true
		case FlagNoIgnoreSigpipe:
			ignoreSigpipe = false
		case FlagIgnoreWriteErrors:
			ignoreWriteErrs = true
		case FlagNoIgnoreWriteErrors:
			ignoreWriteErrs = false
		}
	})

	if ignoreSigpipe {
		signal.Ignore(syscall.SIGPIPE)
	}

	cmds := make([]*exec.Cmd, 0, len(args))
	pipes := make([]io.Writer, 0, len(args))
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
		if ignoreWriteErrs {
			pipes = append(pipes, util.NewSuppressErrorWriter(stdin))
		} else {
			pipes = append(pipes, stdin)
		}
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Wait for all commands to exit
		for _, e := range cmds {
			if err := e.Wait(); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}

		// Close pipes
		for _, pipe := range pipes {
			_ = pipe.(io.Closer).Close()
		}
	}()

	if _, err := io.Copy(io.MultiWriter(pipes...), cmd.InOrStdin()); err != nil {
		mu.Lock()
		if ignoreWriteErrs {
			errs = append(errs, util.NewExitCodeError(1))
		} else {
			errs = append(errs, err)
		}
		mu.Unlock()
	}
	for _, pipe := range pipes {
		_ = pipe.(io.Closer).Close()
	}

	wg.Wait()
	return util.JoinErrors(errs...)
}
