package sponge

import (
	"io"
	"os"
	"path/filepath"

	"gabe565.com/moreutils/internal/cmdutil"
	"gabe565.com/moreutils/internal/util"
	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/must"
	"gabe565.com/utils/termx"
	"github.com/spf13/cobra"
)

const (
	Name       = "sponge"
	FlagAppend = "append"
)

func New(opts ...cobrax.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " file",
		Short:   "Soak up standard input and write to a file",
		RunE:    run,
		GroupID: cmdutil.Applet,
	}

	cmd.Flags().BoolP(FlagAppend, "a", false, "Append to the file instead of overwriting")

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

	var out io.Writer
	var tmp *os.File
	var stat os.FileInfo
	if len(args) == 0 {
		out = cmd.OutOrStdout()
	} else {
		var err error
		stat, err = os.Stat(args[0])
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		tmp, err = os.CreateTemp("", "sponge-*-"+filepath.Base(args[0]))
		if err != nil {
			return err
		}
		defer func() {
			_ = tmp.Close()
			_ = os.Remove(tmp.Name())
		}()

		if stat != nil {
			prevUmask := util.Umask(0)
			if err := tmp.Chmod(stat.Mode()); err != nil {
				return err
			}
			util.Umask(prevUmask)
		}

		out = tmp
	}

	if must.Must2(cmd.Flags().GetBool(FlagAppend)) {
		if in, err := os.Open(args[0]); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			defer func() {
				_ = in.Close()
			}()

			if _, err := io.Copy(out, in); err != nil {
				return err
			}

			_ = in.Close()
		}
	}

	if _, err := io.Copy(out, cmd.InOrStdin()); err != nil {
		return err
	}

	if tmp != nil {
		if err := tmp.Close(); err != nil {
			return err
		}

		if err := os.Rename(tmp.Name(), args[0]); err != nil {
			// Atomic copy not possible
			in, err := os.Open(tmp.Name())
			if err != nil {
				return err
			}
			defer func() {
				_ = in.Close()
			}()

			mode := os.FileMode(0o666)
			if stat != nil {
				mode = stat.Mode()
			}

			out, err := os.OpenFile(args[0], os.O_WRONLY|os.O_TRUNC, mode)
			if err != nil {
				return err
			}

			if _, err := io.Copy(out, in); err != nil {
				return err
			}

			if err := out.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}
