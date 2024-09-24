package errno

import (
	"fmt"
	"io"
	"iter"
	"strconv"
	"strings"

	"github.com/gabe565/moreutils/internal/util"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

const Supported = true

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	if util.Must2(cmd.Flags().GetBool(FlagList)) {
		for errno, name := range iterErrnos() {
			if err := printErrno(cmd.OutOrStdout(), name, errno); err != nil {
				return err
			}
		}
		return nil
	}

	if util.Must2(cmd.Flags().GetBool(FlagSearch)) {
		var searchStr string
		if len(args) != 0 {
			searchStr = strings.ToLower(args[0])
		}
		for errno, name := range iterErrnos() {
			if strings.Contains(errno.Error(), searchStr) {
				if err := printErrno(cmd.OutOrStdout(), name, errno); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if len(args) == 0 {
		return nil
	}

	// Get code by name
	if code, err := strconv.Atoi(args[0]); err == nil {
		if name := unix.ErrnoName(unix.Errno(code)); name != "" {
			return printErrno(cmd.OutOrStdout(), name, unix.Errno(code))
		}
		return fmt.Errorf("%w: %d", ErrUnknownNo, code)
	}

	// Get code by number
	for errno, name := range iterErrnos() {
		if name == args[0] {
			return printErrno(cmd.OutOrStdout(), name, errno)
		}
	}
	return fmt.Errorf("%w: %s", ErrUnknown, args[0])
}

func printErrno(w io.Writer, name string, errno unix.Errno) error {
	_, err := fmt.Fprintf(w, "%s %d %s\n", name, uint(errno), errno.Error())
	return err
}

func iterErrnos() iter.Seq2[unix.Errno, string] {
	return func(yield func(unix.Errno, string) bool) {
		for num := range 256 {
			if name := unix.ErrnoName(unix.Errno(num)); name != "" {
				if !yield(unix.Errno(num), name) {
					return
				}
			}
		}

		// MIPS has an error code > 256
		if name := unix.ErrnoName(unix.Errno(1133)); name != "" {
			if !yield(unix.Errno(1133), name) {
				return
			}
		}
	}
}