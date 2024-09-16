package mispipe

import (
	"os/exec"
	"sync"

	"github.com/gabe565/moreutils/cmd/cmdutil"
	"github.com/spf13/cobra"
)

const Name = "mispipe"

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " command command",
		Short:   "Pipe two commands, returning the exit status of the first",
		Args:    cobra.ExactArgs(2),
		RunE:    run,
		GroupID: cmdutil.Applet,
	}
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	e1 := exec.Command("sh", "-c", args[0])
	e1.Stdin = cmd.InOrStdin()
	e1.Stderr = cmd.ErrOrStderr()

	e2 := exec.Command("sh", "-c", args[1])
	e2.Stdout = cmd.OutOrStdout()
	e2.Stderr = cmd.ErrOrStderr()
	var err error
	e2.Stdin, err = e1.StdoutPipe()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = e2.Run()
	}()

	if err := e1.Start(); err != nil {
		return err
	}

	wg.Wait()
	return e1.Wait()
}
