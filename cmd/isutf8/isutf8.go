package isutf8

import (
	"bufio"
	"errors"
	"fmt"
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
		Args:    cobra.MinimumNArgs(1),
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
	for _, arg := range args {
		if err := checkFile(cmd, arg); err != nil {
			if errors.Is(err, errNotUTF8) {
				notAllUTF8 = true
			} else {
				errs = append(errs, err)
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

	scanner := bufio.NewScanner(f)
	var line, offset int64
	for scanner.Scan() {
		line++
		b := scanner.Bytes()

		for i := 0; i < len(b); {
			if r, width := utf8.DecodeRune(b[i:]); width != 0 {
				if r == utf8.RuneError {
					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s: line %d, char %d, byte %d\n",
						path, line, i, offset)
					return fmt.Errorf("%s: %w", path, errNotUTF8)
				}
				i += width
				offset++
			}
		}
	}
	return scanner.Err()
}
