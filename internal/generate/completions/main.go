package main

import (
	"os"
	"path/filepath"
	"slices"

	"gabe565.com/moreutils/cmd"
	"gabe565.com/moreutils/internal/cmdutil/subcommands"
	flag "github.com/spf13/pflag"
)

const (
	shellBash = "bash"
	shellZsh  = "zsh"
	shellFish = "fish"
)

func main() {
	flags := flag.NewFlagSet("", flag.ContinueOnError)

	if err := flags.Parse(os.Args); err != nil {
		panic(err)
	}

	if err := os.RemoveAll("completions"); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("completions", 0o777); err != nil {
		panic(err)
	}

	cmds := append(slices.Collect(subcommands.Without(nil)), cmd.New(cmd.Name))
	for _, shell := range []string{shellBash, shellZsh, shellFish} {
		basePath := filepath.Join("completions", shell)
		if err := os.MkdirAll(basePath, 0o777); err != nil {
			panic(err)
		}
		for _, subCmd := range cmds {
			var path string
			switch shell {
			case shellBash:
				path = filepath.Join(basePath, subCmd.Name())
				if err := subCmd.GenBashCompletionFileV2(path, true); err != nil {
					panic(err)
				}
			case shellZsh:
				path = filepath.Join(basePath, "_"+subCmd.Name())
				if err := subCmd.GenZshCompletionFile(path); err != nil {
					panic(err)
				}
			case shellFish:
				path = filepath.Join(basePath, subCmd.Name()+".fish")
				if err := subCmd.GenFishCompletionFile(path, true); err != nil {
					panic(err)
				}
			}
		}
	}
}
