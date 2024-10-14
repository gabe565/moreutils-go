package vipe

import (
	"io"
	"os"
	"strings"

	"gabe565.com/moreutils/internal/cmdutil"
	"gabe565.com/moreutils/internal/editor"
	"gabe565.com/moreutils/internal/util"
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
	cmd.SilenceUsage = true

	suffix := util.Must2(cmd.Flags().GetString(FlagSuffix))
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

	if !util.IsTerminal(cmd.InOrStdin()) {
		if _, err := io.Copy(tmp, cmd.InOrStdin()); err != nil {
			return err
		}
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	_, disableTTY := cmd.Annotations[cmdutil.DisableTTYAnnotation]
	if err := editor.Edit(tmp.Name(), !disableTTY); err != nil {
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
