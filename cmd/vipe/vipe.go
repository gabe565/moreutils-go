package vipe

import (
	"io"
	"os"
	"strings"

	"github.com/gabe565/moreutils/cmd/cmdutil"
	"github.com/gabe565/moreutils/internal/editor"
	"github.com/gabe565/moreutils/internal/util"
	"github.com/spf13/cobra"
)

const (
	Name       = "vipe"
	FlagSuffix = "suffix"
)

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name,
		Short:   "Insert a text editor into a pipe",
		Args:    cobra.NoArgs,
		RunE:    run,
		GroupID: cmdutil.Applet,

		ValidArgsFunction: cobra.NoFileCompletions,
	}

	cmd.Flags().StringP(FlagSuffix, "s", "txt", "Provides a file extension to the temp file generated")

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
	if util.IsTerminal(cmd.InOrStdin()) {
		return cmd.Usage()
	}
	cmd.SilenceUsage = true

	suffix, err := cmd.Flags().GetString(FlagSuffix)
	if err != nil {
		panic(err)
	}
	if suffix != "" && !strings.HasPrefix(suffix, ".") {
		suffix = "." + suffix
	}

	tmp, err := os.CreateTemp("", "vipe-*"+suffix)
	if err != nil {
		return err
	}
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
	}()

	if _, err := io.Copy(tmp, cmd.InOrStdin()); err != nil {
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	if err := editor.Edit(tmp.Name(), true); err != nil {
		return err
	}

	tmp, err = os.Open(tmp.Name())
	if err != nil {
		return err
	}
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
	}()

	if _, err := io.Copy(cmd.OutOrStdout(), tmp); err != nil {
		return err
	}

	return nil
}
