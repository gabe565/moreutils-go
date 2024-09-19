package subcommands

import (
	"errors"
	"fmt"
	"iter"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gabe565/moreutils/cmd/chronic"
	"github.com/gabe565/moreutils/cmd/cmdutil"
	"github.com/gabe565/moreutils/cmd/combine"
	"github.com/gabe565/moreutils/cmd/errno"
	"github.com/gabe565/moreutils/cmd/ifne"
	"github.com/gabe565/moreutils/cmd/mispipe"
	"github.com/gabe565/moreutils/cmd/parallel"
	"github.com/gabe565/moreutils/cmd/pee"
	"github.com/gabe565/moreutils/cmd/sponge"
	"github.com/gabe565/moreutils/cmd/ts"
	"github.com/gabe565/moreutils/cmd/vidir"
	"github.com/gabe565/moreutils/cmd/vipe"
	"github.com/gabe565/moreutils/cmd/zrun"
	"github.com/spf13/cobra"
)

func All(opts ...cmdutil.Option) []*cobra.Command {
	cmds := make([]*cobra.Command, 0, 12)
	cmds = append(cmds,
		chronic.New(opts...),
		combine.New(opts...),
	)
	if errno.Supported {
		cmds = append(cmds, errno.New(opts...))
	}
	cmds = append(cmds,
		ifne.New(opts...),
		mispipe.New(opts...),
		parallel.New(opts...),
		pee.New(opts...),
		sponge.New(opts...),
		ts.New(opts...),
		vidir.New(opts...),
		vipe.New(opts...),
		zrun.New(opts...),
	)
	return cmds
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
		if errno.Supported {
			return errno.New(opts...), nil
		}
	case ifne.Name:
		return ifne.New(opts...), nil
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
