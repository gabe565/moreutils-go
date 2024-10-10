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

	if len(args) == 0 {
		cmd.SilenceUsage = false
		return cobra.MinimumNArgs(1)(cmd, args)
	}

	if util.Must2(cmd.Flags().GetBool(FlagSearch)) {
		search := strings.ToLower(strings.Join(args, " "))
		for e := range errno.Iter() {
			if strings.Contains(e.Error(), search) {
				if err := printErrno(cmd, e); err != nil {
					return err
				}
			}
		}
		return nil
	}

	var errs []error
	for _, arg := range args {
		if err := findErrno(cmd, arg); err != nil {
			errs = append(errs, err)
		}
	}
	return util.JoinErrors(errs...)
}

func findErrno(cmd *cobra.Command, arg string) error {
	// Get code by name
	if code, err := strconv.Atoi(arg); err == nil {
		if e := errno.New(code); e.Valid() {
			return printErrno(cmd, e)
		}
		return fmt.Errorf("%w: %d", ErrUnknownNo, code)
	}

	// Get code by number
	for e := range errno.Iter() {
		if e.Name() == arg {
			return printErrno(cmd, e)
		}
	}
	return fmt.Errorf("%w: %s", ErrUnknown, arg)
}

func printErrno(cmd *cobra.Command, e *errno.Errno) error {
	pretty := e.Error()
	pretty = strings.ToUpper(pretty[0:1]) + pretty[1:]
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %d %s\n", e.Name(), uint(e.Errno), pretty)
	return err
}
