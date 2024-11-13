package ifne

import (
	"bufio"
	"errors"
	"io"
	"os/exec"

	"gabe565.com/moreutils/internal/cmdutil"
	"gabe565.com/moreutils/internal/util"
	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/must"
	"gabe565.com/utils/termx"
	"github.com/spf13/cobra"
)

const (
	Name       = "ifne"
	FlagInvert = "invert"
)

func New(opts ...cobrax.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " [flags] command",
		Short:   "Run a command if the standard input is not empty",
		Args:    cobra.MinimumNArgs(1),
		RunE:    run,
		GroupID: cmdutil.Applet,

		DisableFlagsInUseLine: true,
	}

	cmd.Flags().SetInterspersed(false)
	cmd.Flags().BoolP(FlagInvert, "n", false, "Inverse operation. Run the command if the standard input is empty.")

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if termx.IsTerminal(cmd.InOrStdin()) {
		return util.ErrNotAPipe
	}
	cmd.SilenceUsage = true

	invert := must.Must2(cmd.Flags().GetBool(FlagInvert))

	r := bufio.NewReader(cmd.InOrStdin())
	willRun, err := shouldRun(r, invert)
	if err != nil {
		return err
	}

	if willRun {
		e := exec.Command(args[0], args[1:]...)
		e.Stdin = r
		e.Stdout = cmd.OutOrStdout()
		e.Stderr = cmd.ErrOrStderr()
		return e.Run()
	}
	return nil
}

func shouldRun(r *bufio.Reader, invert bool) (bool, error) {
	if _, err := r.ReadByte(); err == nil {
		if err := r.UnreadByte(); err != nil {
			return false, err
		}
	} else {
		if errors.Is(err, io.EOF) {
			return invert, nil
		}
		return false, err
	}
	return !invert, nil
}
