package main

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"gabe565.com/moreutils/cmd"
	"gabe565.com/moreutils/cmd/ifdata"
	"gabe565.com/moreutils/internal/cmdutil/subcommands"
	"gabe565.com/moreutils/internal/generate/seealsoreplacer"
	"gabe565.com/moreutils/internal/util"
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra/doc"
	flag "github.com/spf13/pflag"
)

func main() {
	flags := flag.NewFlagSet("", flag.ContinueOnError)

	var version string
	flags.StringVar(&version, "version", "beta", "Version")

	var dateParam string
	flags.StringVar(&dateParam, "date", time.Now().Format(time.RFC3339), "Build date")

	if err := flags.Parse(os.Args); err != nil {
		panic(err)
	}

	if err := os.RemoveAll("manpages"); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("manpages", 0o755); err != nil {
		panic(err)
	}

	date, err := time.Parse(time.RFC3339, dateParam)
	if err != nil {
		panic(err)
	}

	opts := []cobrax.Option{cobrax.WithVersion("beta")}
	root := cmd.New(cmd.Name, opts...)
	cmds := append(slices.Collect(subcommands.Without(nil, opts...)), root)
	for _, subCmd := range root.Commands() {
		// Add any commands which aren't standalone
		if !util.CmdsContains(cmds, subCmd) {
			cmds = append(cmds, subCmd)
		}
	}

	linked := slices.Collect(subcommands.Without(nil))
	for _, subCmd := range cmds {
		subCmd.DisableAutoGenTag = true
		header := doc.GenManHeader{
			Title:   strings.ToUpper(subCmd.Name()),
			Section: "1",
			Date:    &date,
			Source:  cmd.Name + " " + version,
			Manual:  "User Commands",
		}

		name := subCmd.Name() + ".1.gz"
		if subCmd.Name() != cmd.Name && !util.CmdsContains(linked, subCmd) {
			name = cmd.Name + "-" + name
		}
		path := filepath.Join("manpages", name)
		out, err := os.Create(path)
		if err != nil {
			panic(err)
		}

		gzw := gzip.NewWriter(out)
		w := io.Writer(gzw)
		if subCmd.Name() == cmd.Name {
			// Replace "See Also" section in root command
			w = seealsoreplacer.New(w, ".SH SEE ALSO\n", linked)
		}

		if subCmd.Name() == ifdata.Name {
			w = optionsInserter{
				w:      w,
				find:   "help for " + ifdata.Name + "\n",
				insert: ifdata.ManOptions(),
			}
		}

		if err := doc.GenMan(subCmd, &header, w); err != nil {
			panic(err)
		}

		if subCmd.Name() != cmd.Name && !subCmd.HasParent() {
			// Add "See Also" section to standalone commands
			_, _ = io.WriteString(w, "\n\n.SH SEE ALSO\n.PP\n\\fB"+cmd.Name+"(1)\\fP\n")
		}

		if err := errors.Join(gzw.Close(), out.Close()); err != nil {
			panic(err)
		}
	}
}

type optionsInserter struct {
	w      io.Writer
	find   string
	insert string
}

func (o optionsInserter) Write(p []byte) (int, error) {
	if !strings.HasSuffix(o.insert, "\n") {
		o.insert += "\n"
	}

	if bytes.Contains(p, []byte(o.find)) {
		beforeIdx := bytes.Index(p, []byte(o.find))
		if beforeIdx == -1 {
			panic("missing header: " + o.find)
		}
		beforeIdx += len(o.find)

		_, err := o.w.Write(slices.Concat(p[:beforeIdx], []byte(o.find), []byte(o.insert), p[beforeIdx:]))
		return len(p), err
	}
	return o.w.Write(p)
}
