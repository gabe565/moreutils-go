package main

import (
	"os"
	"path/filepath"

	"gabe565.com/moreutils/cmd"
	"gabe565.com/moreutils/internal/cmdutil/subcommands"
)

func main() {
	if err := os.RemoveAll("links"); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("links", 0o755); err != nil {
		panic(err)
	}

	for subCmd := range subcommands.Without(nil) {
		path := filepath.Join("links", subCmd.Name())
		if err := os.Symlink(cmd.Name, path); err != nil {
			panic(err)
		}
	}
}
