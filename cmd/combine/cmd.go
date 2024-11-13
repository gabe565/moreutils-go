package combine

import (
	"bytes"
	"errors"
	"io"
	"os"

	"gabe565.com/moreutils/internal/cmdutil"
	"gabe565.com/moreutils/internal/seekbuf"
	"gabe565.com/moreutils/internal/util"
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra"
)

const (
	Name  = "combine"
	Alias = "_"
)

func New(opts ...cobrax.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " file operator file",
		Short:   "Combine sets of lines from two files using boolean operations",
		Args:    cobra.RangeArgs(3, 4),
		RunE:    run,
		GroupID: cmdutil.Applet,

		ValidArgsFunction: validArgs,
	}
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func validArgs(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 1 {
		return operatorStrings(), cobra.ShellCompDirectiveNoFileComp
	}
	return nil, cobra.ShellCompDirectiveDefault
}

var ErrBothFilesStdin = errors.New("both files are stdin")

func run(cmd *cobra.Command, args []string) error {
	if (args[0] == "-" || args[2] == "-") && util.IsTerminal(cmd.InOrStdin()) {
		return util.ErrNotAPipe
	}

	if args[0] == "-" && args[2] == "-" {
		return ErrBothFilesStdin
	}

	op, err := operatorString(args[1])
	if err != nil {
		return err
	}

	f1, err := openFile(cmd, args[0])
	if err != nil {
		return err
	}
	defer func() {
		_ = f1.Close()
	}()

	f2, err := openFile(cmd, args[2])
	if err != nil {
		return err
	}
	defer func() {
		_ = f2.Close()
	}()

	cmd.SilenceUsage = true
	return op.compare(cmd.OutOrStdout(), f1, f2)
}

// openFile opens the given file, or buffers stdin if "-"
func openFile(cmd *cobra.Command, path string) (io.ReadSeekCloser, error) {
	switch path {
	case "-":
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, cmd.InOrStdin()); err != nil {
			return nil, err
		}
		return seekbuf.New(buf.Bytes()), nil
	default:
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
}
