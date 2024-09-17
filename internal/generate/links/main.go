package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/gabe565/moreutils/cmd"
	"github.com/gabe565/moreutils/cmd/cmdutil/subcommands"
	flag "github.com/spf13/pflag"
	"golang.org/x/sys/unix"
)

func main() {
	flags := flag.NewFlagSet("", flag.ContinueOnError)

	var dateParam int64
	flags.Int64Var(&dateParam, "date", time.Now().Unix(), "Build date")

	if err := flags.Parse(os.Args); err != nil {
		panic(err)
	}

	date := unix.NsecToTimeval(time.Unix(dateParam, 0).UnixNano())

	if err := os.RemoveAll("links"); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("links", 0o755); err != nil {
		panic(err)
	}

	for _, subCmd := range subcommands.All() {
		path := filepath.Join("links", subCmd.Name())
		if err := os.Symlink(cmd.Name, path); err != nil {
			panic(err)
		}

		if err := unix.Lutimes(path, []unix.Timeval{date, date}); err != nil {
			panic(err)
		}
	}
}
