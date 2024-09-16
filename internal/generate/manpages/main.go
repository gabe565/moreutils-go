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

	"github.com/gabe565/moreutils/cmd"
	"github.com/gabe565/moreutils/cmd/cmdutil/subcommands"
	"github.com/spf13/cobra"
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

	root := cmd.New(cmd.Name)
	cmds := append(subcommands.All(), root)
	for _, subCmd := range root.Commands() {
		// Add any commands which aren't standalone
		if !slices.ContainsFunc(cmds, func(cmd *cobra.Command) bool { return cmd.Name() == subCmd.Name() }) {
			cmds = append(cmds, subCmd)
		}
	}

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
		if subCmd.HasParent() {
			name = subCmd.Parent().Name() + "-" + name
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
			w = &seeAlsoWriter{w: gzw}
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

type seeAlsoWriter struct {
	w       io.Writer
	replace bool
}

func (s *seeAlsoWriter) Write(p []byte) (int, error) {
	const seeAlsoLine = ".SH SEE ALSO"
	if bytes.Contains(p, []byte(seeAlsoLine)) {
		s.replace = true
		n := len(p)
		headerIdx := bytes.Index(p, []byte(seeAlsoLine))
		p, temp := p[:headerIdx], p[headerIdx:]
		for _, subCmd := range subcommands.All() {
			temp = bytes.ReplaceAll(temp, []byte(cmd.Name+"-"+subCmd.Name()), []byte(subCmd.Name()))
		}
		_, err := s.w.Write(append(p, temp...))
		return n, err
	}
	return s.w.Write(p)
}
