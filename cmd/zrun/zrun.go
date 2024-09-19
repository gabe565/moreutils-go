package zrun

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gabe565/moreutils/cmd/cmdutil"
	"github.com/spf13/cobra"
)

const (
	Name   = "zrun"
	Prefix = "z"
)

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " command args",
		Short:   "Automatically decompress arguments to command",
		RunE:    run,
		GroupID: cmdutil.Applet,

		DisableFlagsInUseLine: true,
	}
	cmd.Flags().SetInterspersed(false)
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	base := filepath.Base(os.Args[0])
	if cmd.Parent() != nil || base == cmd.Name() || !strings.HasPrefix(base, "z") {
		// Default behavior
		if err := cobra.MinimumNArgs(2)(cmd, args); err != nil {
			return err
		}
	} else {
		// Command is linked to another command
		args = append([]string{strings.TrimPrefix(base, "z")}, args...)
	}

	cmd.SilenceUsage = true

	for i, arg := range args {
		if i == 0 {
			// Do not mutate the command
			continue
		}

		if isCompressed(arg) {
			arg, err := decompressTmp(cmd, arg)
			defer func() {
				if arg != "" {
					_ = os.Remove(arg)
				}
			}()
			if err != nil {
				return err
			}

			args[i] = arg
		}
	}

	e := exec.Command(args[0], args[1:]...)
	e.Stdin = cmd.InOrStdin()
	e.Stdout = cmd.OutOrStdout()
	e.Stderr = cmd.ErrOrStderr()
	return e.Run()
}

func isCompressed(arg string) bool {
	ext := filepath.Ext(arg)
	return strings.EqualFold(ext, ".gz") || ext == ".Z" ||
		strings.EqualFold(ext, ".bz2") ||
		strings.EqualFold(ext, ".xz") ||
		strings.EqualFold(ext, ".lzma") ||
		strings.EqualFold(ext, ".lzo")
}

var ErrUnknownExtension = errors.New("unknown extension")

func decompressTmp(cmd *cobra.Command, path string) (string, error) {
	in, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = in.Close()
	}()

	ext := filepath.Ext(path)
	withoutExt := strings.TrimSuffix(path, ext)
	tmp, err := os.CreateTemp("", "zrun-*-"+filepath.Base(withoutExt))
	if err != nil {
		return "", err
	}
	defer func() {
		_ = tmp.Close()
	}()

	switch {
	case strings.EqualFold(ext, ".gz"), ext == ".Z":
		gzr, err := gzip.NewReader(in)
		if err != nil {
			return tmp.Name(), err
		}

		if _, err := io.Copy(tmp, gzr); err != nil {
			return tmp.Name(), err
		}

		if err := gzr.Close(); err != nil {
			return tmp.Name(), err
		}
	default:
		var args []string
		switch {
		case strings.EqualFold(ext, ".bz2"):
			args = []string{"bzip2", "-d", "-c"}
		case strings.EqualFold(ext, ".xz"):
			args = []string{"xz", "-d", "-c"}
		case strings.EqualFold(ext, ".lzma"):
			args = []string{"lzma", "-d", "-c"}
		case strings.EqualFold(ext, ".lzo"):
			args = []string{"lzop", "-d", "-c"}
		default:
			return tmp.Name(), fmt.Errorf("%w: %s", ErrUnknownExtension, ext)
		}

		if err := execDecompress(args, in, tmp, cmd.ErrOrStderr()); err != nil {
			return tmp.Name(), err
		}
	}

	if err := tmp.Close(); err != nil {
		return tmp.Name(), err
	}

	return tmp.Name(), nil
}

func execDecompress(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}
