package subcommands

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gabe565/moreutils/cmd/chronic"
	"github.com/gabe565/moreutils/cmd/cmdutil"
	"github.com/gabe565/moreutils/cmd/combine"
	"github.com/gabe565/moreutils/cmd/ifne"
	"github.com/gabe565/moreutils/cmd/mispipe"
	"github.com/gabe565/moreutils/cmd/pee"
	"github.com/gabe565/moreutils/cmd/sponge"
	"github.com/gabe565/moreutils/cmd/ts"
	"github.com/gabe565/moreutils/cmd/vidir"
	"github.com/gabe565/moreutils/cmd/vipe"
	"github.com/gabe565/moreutils/cmd/zrun"
	"github.com/spf13/cobra"
)

func All(opts ...cmdutil.Option) []*cobra.Command {
	return []*cobra.Command{
		chronic.New(opts...),
		combine.New(opts...),
		ifne.New(opts...),
		mispipe.New(opts...),
		pee.New(opts...),
		sponge.New(opts...),
		ts.New(opts...),
		vidir.New(opts...),
		vipe.New(opts...),
		zrun.New(opts...),
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
	case ifne.Name:
		return ifne.New(opts...), nil
	case mispipe.Name:
		return mispipe.New(opts...), nil
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
	default:
		if strings.HasPrefix(base, zrun.Prefix) {
			return zrun.New(opts...), nil
		}

		return nil, fmt.Errorf("%w: %s", ErrUnknownCommand, base)
	}
}
