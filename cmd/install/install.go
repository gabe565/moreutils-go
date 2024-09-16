package install

import (
	"os"
	"path"

	"github.com/gabe565/moreutils/cmd/cmdutil/subcommands"
	"github.com/gabe565/moreutils/internal/util"
	"github.com/spf13/cobra"
)

const (
	FlagSymbolic = "symbolic"
	FlagForce    = "force"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install dir",
		Short: "Creates hardlinks/symlinks for each applet",
		Args:  cobra.MinimumNArgs(1),
		RunE:  run,
	}

	cmd.Flags().BoolP(FlagSymbolic, "s", false, "Create symbolic links instead of hard links")
	cmd.Flags().BoolP(FlagForce, "f", false, "Overwrite existing files")

	return cmd
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

	force, err := cmd.Flags().GetBool(FlagForce)
	if err != nil {
		panic(err)
	}

	symbolic, err := cmd.Flags().GetBool(FlagSymbolic)
	if err != nil {
		panic(err)
	}

	var errs []error
	for _, cmd := range subcommands.All() {
		dst := path.Join(dst, cmd.Name())
		if err := link(symbolic, src, dst); err != nil {
			if force {
				if err := os.Remove(dst); err == nil {
					if err := link(symbolic, src, dst); err != nil {
						errs = append(errs, err)
					}
				} else {
					errs = append(errs, err)
				}
			} else {
				errs = append(errs, err)
			}
		}
	}

	return util.JoinErrors(errs...)
}

func link(symbolic bool, oldname, newname string) error {
	if symbolic {
		return os.Symlink(oldname, newname)
	}
	return os.Link(oldname, newname)
}
