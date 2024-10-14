package subcommands

import (
	"errors"
	"fmt"
	"iter"
	"path/filepath"
	"slices"
	"strings"

	"gabe565.com/moreutils/cmd/chronic"
	"gabe565.com/moreutils/cmd/combine"
	"gabe565.com/moreutils/cmd/errno"
	"gabe565.com/moreutils/cmd/ifdata"
	"gabe565.com/moreutils/cmd/ifne"
	"gabe565.com/moreutils/cmd/isutf8"
	"gabe565.com/moreutils/cmd/mispipe"
	"gabe565.com/moreutils/cmd/parallel"
	"gabe565.com/moreutils/cmd/pee"
	"gabe565.com/moreutils/cmd/sponge"
	"gabe565.com/moreutils/cmd/ts"
	"gabe565.com/moreutils/cmd/vidir"
	"gabe565.com/moreutils/cmd/vipe"
	"gabe565.com/moreutils/cmd/zrun"
	"gabe565.com/moreutils/internal/cmdutil"
	"github.com/spf13/cobra"
)

func All(opts ...cmdutil.Option) []*cobra.Command {
	return []*cobra.Command{
		chronic.New(opts...),
		combine.New(opts...),
		errno.New(opts...),
		ifdata.New(opts...),
		ifne.New(opts...),
		isutf8.New(opts...),
		mispipe.New(opts...),
		parallel.New(opts...),
		pee.New(opts...),
		sponge.New(opts...),
		ts.New(opts...),
		vidir.New(opts...),
		vipe.New(opts...),
		zrun.New(opts...),
	}
}

func DefaultExcludes() []string {
	return []string{parallel.Name}
}

func Without(excludes []string, opts ...cmdutil.Option) iter.Seq[*cobra.Command] {
	if len(excludes) == 0 {
		excludes = DefaultExcludes()
	}
	return func(yield func(*cobra.Command) bool) {
		for _, cmd := range All(opts...) {
			if !slices.Contains(excludes, cmd.Name()) {
				if !yield(cmd) {
					return
				}
			}
		}
	}
}

var ErrUnknownCommand = errors.New("unknown command")

func Choose(name string, opts ...cmdutil.Option) (*cobra.Command, error) {
	base := filepath.Base(name)
	switch base {
	case chronic.Name:
		return chronic.New(opts...), nil
	case combine.Name, combine.Alias:
		return combine.New(opts...), nil
	case errno.Name:
		return errno.New(opts...), nil
	case ifdata.Name:
		return ifdata.New(opts...), nil
	case ifne.Name:
		return ifne.New(opts...), nil
	case isutf8.Name:
		return isutf8.New(opts...), nil
	case mispipe.Name:
		return mispipe.New(opts...), nil
	case parallel.Name:
		return parallel.New(opts...), nil
	case pee.Name:
		return pee.New(opts...), nil
	case sponge.Name:
		return sponge.New(opts...), nil
	case ts.Name:
		return ts.New(opts...), nil
	case vidir.Name:
		return vidir.New(opts...), nil
	case vipe.Name:
		return vipe.New(opts...), nil
	case zrun.Name:
		return zrun.New(opts...), nil
	}

	if strings.HasPrefix(base, zrun.Prefix) {
		return zrun.New(opts...), nil
	}

	return nil, fmt.Errorf("%w: %s", ErrUnknownCommand, base)
}
