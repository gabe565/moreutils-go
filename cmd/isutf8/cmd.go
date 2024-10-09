package isutf8

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"github.com/gabe565/moreutils/internal/cmdutil"
	"github.com/gabe565/moreutils/internal/util"
	"github.com/spf13/cobra"
)

const (
	Name = "isutf8"
)

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " file...",
		Short:   "Check whether files are valid UTF-8",
		RunE:    run,
		GroupID: cmdutil.Applet,

		DisableFlagsInUseLine: true,
	}

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

var ErrNotAllUTF8 = errors.New("not all UTF-8 encoded")

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	var notAllUTF8 bool
	var errs []error
	if len(args) == 0 {
		if util.IsTerminal(cmd.InOrStdin()) {
			return cmd.Help()
		}

		if err := checkReader(cmd, "(standard input)", cmd.InOrStdin()); err != nil {
			if errors.Is(err, errNotUTF8) {
				notAllUTF8 = true
			} else {
				errs = append(errs, err)
			}
		}
	} else {
		for _, arg := range args {
			if err := checkFile(cmd, arg); err != nil {
				if errors.Is(err, errNotUTF8) {
					notAllUTF8 = true
				} else {
					errs = append(errs, err)
				}
			}
		}
	}

	if notAllUTF8 {
		errs = append(errs, ErrNotAllUTF8)
	}
	return util.JoinErrors(errs...)
}

var errNotUTF8 = errors.New("file is not UTF-8 encoded")

func checkFile(cmd *cobra.Command, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	return checkReader(cmd, path, f)
}

func checkReader(cmd *cobra.Command, path string, r io.Reader) error {
	buf := bufio.NewReader(r)
	line := int64(1)
	var i, offset int64
	for {
		r, size, err := buf.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		switch r {
		case '\n':
			i = 0
			line++
		case utf8.RuneError:
			if _, err := fmt.Fprintf(cmd.OutOrStdout(),
				"%s: line %d, char %d, byte %d\n",
				path, line, i, offset,
			); err != nil {
				return err
			}
			return fmt.Errorf("%s: %w", path, errNotUTF8)
		}

		i += int64(size)
		offset++
	}
}
