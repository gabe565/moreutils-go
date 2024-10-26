package install

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"gabe565.com/moreutils/internal/cmdutil/subcommands"
	"gabe565.com/utils/must"
	"github.com/spf13/cobra"
)

const (
	FlagSymbolic = "symbolic"
	FlagForce    = "force"
	FlagRelative = "relative"
	FlagExclude  = "exclude"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install dir",
		Short: "Creates hardlinks/symlinks for each applet",
		Args:  cobra.ExactArgs(1),
		RunE:  run,

		ValidArgsFunction: validArgs,
	}

	cmd.Flags().BoolP(FlagSymbolic, "s", false, "Create symbolic links instead of hard links")
	cmd.Flags().BoolP(FlagForce, "f", false, "Overwrite existing files")
	cmd.Flags().BoolP(FlagRelative, "r", false, "Create relative symbolic links")
	cmd.Flags().StringSlice(FlagExclude, subcommands.DefaultExcludes(), "Subcommands that will not be linked")

	return cmd
}

func validArgs(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return nil, cobra.ShellCompDirectiveFilterDirs
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	src, err := os.Executable()
	if err != nil {
		return err
	}

	dst := args[0]
	if err := os.MkdirAll(dst, 0o777); err != nil {
		return err
	}

	force := must.Must2(cmd.Flags().GetBool(FlagForce))
	symbolic := must.Must2(cmd.Flags().GetBool(FlagSymbolic))

	if must.Must2(cmd.Flags().GetBool(FlagRelative)) {
		dstAbs, err := filepath.Abs(dst)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dstAbs, filepath.Dir(src))
		if err != nil {
			return err
		}

		src = filepath.Join(relPath, filepath.Base(src))
	}

	excludes := must.Must2(cmd.Flags().GetStringSlice(FlagExclude))

	var errs []error
	for subCmd := range subcommands.Without(excludes) {
		dst := path.Join(dst, subCmd.Name())
		if force {
			if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
				errs = append(errs, err)
				continue
			}
		}
		if err := link(symbolic, src, dst); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func link(symbolic bool, oldname, newname string) error {
	if symbolic {
		return os.Symlink(oldname, newname)
	}
	return os.Link(oldname, newname)
}
