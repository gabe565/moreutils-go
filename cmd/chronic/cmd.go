package chronic

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"

	"gabe565.com/moreutils/internal/cmdutil"
	"gabe565.com/moreutils/internal/execbuf"
	"gabe565.com/utils/must"
	"github.com/spf13/cobra"
)

const (
	Name        = "chronic"
	FlagStderr  = "stderr"
	FlagVerbose = "verbose"
)

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " [flags] command",
		Short:   "Runs a command quietly unless it fails",
		Args:    cobra.MinimumNArgs(1),
		RunE:    run,
		GroupID: cmdutil.Applet,

		DisableFlagsInUseLine: true,
	}

	cmd.Flags().SetInterspersed(false)
	cmd.Flags().BoolP(FlagStderr, "e", false, "Triggers output when stderr output length is non-zero")
	cmd.Flags().BoolP(FlagVerbose, "v", false, "Verbose output (distinguishes between STDOUT and STDERR, also reports RETVAL)")

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	onStderr := must.Must2(cmd.Flags().GetBool(FlagStderr))
	verbose := must.Must2(cmd.Flags().GetBool(FlagVerbose))

	e := exec.Command(args[0], args[1:]...)
	e.Stdin = cmd.InOrStdin()

	buf, err := execbuf.RunBuffered(e, cmd.OutOrStdout(), cmd.ErrOrStderr())
	if err != nil {
		var execErr *exec.ExitError
		if errors.As(err, &execErr) {
			if printErr := printBuf(cmd, buf, execErr.ExitCode(), verbose); printErr != nil {
				err = errors.Join(err, printErr)
			}
		}
		return err
	}

	if onStderr && buf.Len(cmd.ErrOrStderr()) != 0 {
		return printBuf(cmd, buf, 0, verbose)
	}

	return nil
}

func printBuf(cmd *cobra.Command, buf *execbuf.Buffer, exitCode int, verbose bool) error {
	if !verbose {
		return buf.Print(nil)
	}

	var errs []error
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "STDOUT:")
	errs = append(errs, buf.Print(cmd.OutOrStdout()))
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "\nSTDERR:")
	errs = append(errs, buf.Print(cmd.ErrOrStderr()))

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "\nRETVAL:", strconv.Itoa(exitCode))
	return errors.Join(errs...)
}
