//go:build unix

package errno

import (
	"fmt"
	"strconv"
	"strings"

	"gabe565.com/moreutils/internal/errno"
	"gabe565.com/moreutils/internal/util"
	"github.com/spf13/cobra"
)

const Supported = true

func validArgs(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	if util.Must2(cmd.Flags().GetBool(FlagList)) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	search := util.Must2(cmd.Flags().GetBool(FlagSearch))

	var completions []string
	for e := range errno.Iter() {
		if search {
			completions = append(completions, e.Error()+"\t"+e.Name()+" "+strconv.Itoa(int(e.Errno)))
		} else {
			num := strconv.Itoa(int(e.Errno))
			completions = append(completions,
				num+"\t"+e.Name()+" "+e.Error(),
				e.Name()+"\t"+num+" "+e.Error(),
			)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

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
			if strings.Contains(strings.ToLower(e.Error()), search) {
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
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %d %s\n", e.Name(), uint(e.Errno), e.Error())
	return err
}
