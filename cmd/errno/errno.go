package errno

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"strconv"
	"strings"

	"github.com/gabe565/moreutils/cmd/cmdutil"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

const (
	Name       = "errno"
	FlagList   = "list"
	FlagSearch = "search"
)

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " { name | code }",
		Short:   "Look up errno names and descriptions",
		Args:    cobra.MaximumNArgs(1),
		RunE:    run,
		GroupID: cmdutil.Applet,
	}

	cmd.Flags().BoolP(FlagList, "l", false, "List all errno values")
	cmd.Flags().BoolP(FlagSearch, "s", false, "Search for errors whose description contains all the given words (case-insensitive)")

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

var (
	ErrUnknownNo = errors.New("unknown errno")
	ErrUnknown   = errors.New("unknown err name")
)

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	if list, err := cmd.Flags().GetBool(FlagList); err != nil {
		panic(err)
	} else if list {
		for errno, name := range iterErrnos() {
			if err := printErrno(cmd.OutOrStdout(), name, errno); err != nil {
				return err
			}
		}
		return nil
	}

	if search, err := cmd.Flags().GetBool(FlagSearch); err != nil {
		panic(err)
	} else if search {
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
