package combine

import (
	"bytes"
	"io"
	"os"

	"github.com/gabe565/moreutils/cmd/cmdutil"
	"github.com/gabe565/moreutils/internal/seekbuf"
	"github.com/gabe565/moreutils/internal/util"
	"github.com/spf13/cobra"
)

const (
	Name  = "combine"
	Alias = "_"
)

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " file operator file",
		Short:   "Combine sets of lines from two files using boolean operations",
		Args:    cobra.RangeArgs(3, 4),
		RunE:    run,
		GroupID: cmdutil.Applet,
	}
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if args[0] == "-" || args[2] == "-" {
		if util.IsTerminal(cmd.InOrStdin()) {
			return cmd.Usage()
		}
	}

	var op operator
	if err := op.UnmarshalText([]byte(args[1])); err != nil {
		return err
	}

	f1, closeF1, err := openFile(cmd, args[0])
	if err != nil {
		return err
	}
	defer closeF1()

	f2, closeF2, err := openFile(cmd, args[2])
	if err != nil {
		return err
	}
	defer closeF2()

	cmd.SilenceUsage = true
	return op.compare(cmd.OutOrStdout(), f1, f2)
}

// openFile opens the given file, or buffers stdin if "-"
func openFile(cmd *cobra.Command, path string) (io.ReadSeeker, func(), error) {
	switch path {
	case "-":
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, cmd.InOrStdin()); err != nil {
			return nil, nil, err
		}
		return seekbuf.New(buf.Bytes()), func() {}, nil
	default:
		f, err := os.Open(path)
		if err != nil {
			return nil, nil, err
		}
		return f, func() { _ = f.Close() }, nil
	}
}
