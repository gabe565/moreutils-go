//go:build unix

package errno

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gabe565/moreutils/internal/errno"
	"github.com/gabe565/moreutils/internal/util"
	"github.com/spf13/cobra"
)

const Supported = true

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	if util.Must2(cmd.Flags().GetBool(FlagList)) {
		for e := range errno.Iter() {
			if err := printErrno(cmd, e); err != nil {
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
		for e := range errno.Iter() {
			if strings.Contains(e.Error(), searchStr) {
				if err := printErrno(cmd, e); err != nil {
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
		if e := errno.New(code); e.Valid() {
			return printErrno(cmd, e)
		}
		return fmt.Errorf("%w: %d", ErrUnknownNo, code)
	}

	// Get code by number
	for e := range errno.Iter() {
		if e.Name() == args[0] {
			return printErrno(cmd, e)
		}
	}
	return fmt.Errorf("%w: %s", ErrUnknown, args[0])
}

func printErrno(cmd *cobra.Command, e *errno.Errno) error {
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %d %s\n", e.Name(), uint(e.Errno), e.Error())
	return err
}
