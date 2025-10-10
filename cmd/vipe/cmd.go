package vipe

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"gabe565.com/moreutils/internal/cmdutil"
	"gabe565.com/moreutils/internal/editor"
	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/must"
	"gabe565.com/utils/termx"
	"github.com/spf13/cobra"
)

const (
	Name       = "vipe"
	FlagSuffix = "suffix"
)

func New(opts ...cobrax.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " [file]",
		Short:   "Insert a text editor into a pipe",
		Args:    cobra.MaximumNArgs(1),
		RunE:    run,
		GroupID: cmdutil.Applet,
	}

	cmd.Flags().StringP(FlagSuffix, "s", "txt", "File extension to use for the temp buffer file")

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	suffix := must.Must2(cmd.Flags().GetString(FlagSuffix))

	if len(args) != 0 && !cmd.Flags().Changed(FlagSuffix) {
		if v := filepath.Ext(args[0]); v != "" {
			suffix = v
		}
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

	switch {
	case len(args) != 0:
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}

		_, err = io.Copy(tmp, f)
		_ = f.Close()

		if err != nil {
			return err
		}
	case !termx.IsTerminal(cmd.InOrStdin()):
		if _, err := io.Copy(tmp, cmd.InOrStdin()); err != nil {
			return err
		}
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	_, disableTTY := cmd.Annotations[cmdutil.DisableTTYAnnotation]
	if err := editor.Edit(cmd.Context(), tmp.Name(), !disableTTY); err != nil {
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
