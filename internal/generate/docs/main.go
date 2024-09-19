package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"

	"github.com/gabe565/moreutils/cmd"
	"github.com/gabe565/moreutils/cmd/cmdutil"
	"github.com/gabe565/moreutils/cmd/cmdutil/subcommands"
	"github.com/gabe565/moreutils/internal/util"
	"github.com/spf13/cobra/doc"
)

var (
	ErrMissingBeginApplets = errors.New("readme missing begin applets comment")
	ErrMissingEndApplets   = errors.New("readme missing end applets comment")
)

func main() {
	output := "./docs"

	if err := os.RemoveAll(output); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(output, 0o755); err != nil {
		panic(err)
	}

	opts := []cmdutil.Option{cmdutil.WithVersion("beta")}
	root := cmd.New(cmd.Name, opts...)
	cmds := append(slices.Collect(subcommands.Without(nil, opts...)), root)
	for _, subCmd := range root.Commands() {
		// Add any commands which aren't standalone
		if !util.CmdsContains(cmds, subCmd) {
			cmds = append(cmds, subCmd)
		}
	}

	for _, subCmd := range cmds {
		subCmd.DisableAutoGenTag = true

		name := subCmd.Name() + ".md"
		if subCmd.HasParent() {
			name = subCmd.Parent().Name() + "_" + name
		}
		path := filepath.Join("docs", name)

		out, err := os.Create(path)
		if err != nil {
			panic(err)
		}

		w := io.Writer(out)
		if subCmd.Name() == cmd.Name {
			// Replace "See Also" section in root command
			w = &seeAlsoWriter{w: w}
		}

		if err := doc.GenMarkdown(subCmd, w); err != nil {
			panic(err)
		}

		if subCmd.Name() != cmd.Name && !subCmd.HasParent() {
			// Add "See Also" section to standalone commands
			_, _ = io.WriteString(w, "### SEE ALSO\n\n* ["+cmd.Name+"]("+cmd.Name+".md)\t - "+root.Short+"\n\n")
		}

		if err := out.Close(); err != nil {
			panic(err)
		}
	}

	readmeContents, err := os.ReadFile("README.md")
	if err != nil {
		panic(err)
	}

	const beforeMarker, afterMarker = "## Applets\n\n", "\n## Installation\n"

	beforeApplets, _, found := bytes.Cut(readmeContents, []byte(beforeMarker))
	if !found {
		panic(ErrMissingBeginApplets)
	}
	_, afterApplets, found := bytes.Cut(readmeContents[len(beforeApplets):], []byte(afterMarker))
	if !found {
		panic(ErrMissingEndApplets)
	}

	var list []byte
	linked := slices.Collect(subcommands.Without(nil))
	for _, subCmd := range subcommands.All() {
		docPath := subCmd.Name() + ".md"
		if !util.CmdsContains(linked, subCmd) {
			docPath = cmd.Name + "_" + docPath
		}
		list = append(list, []byte("- **["+subCmd.Name()+"](docs/"+docPath+")**: "+subCmd.Short+"\n")...)
	}

	readmeContents = slices.Concat(beforeApplets, []byte(beforeMarker), list, []byte(afterMarker), afterApplets)

	//nolint:gosec
	if err := os.WriteFile("README.md", readmeContents, 0o644); err != nil {
		panic(err)
	}
}

type seeAlsoWriter struct {
	w       io.Writer
	replace bool
}

func (s *seeAlsoWriter) Write(p []byte) (int, error) {
	const seeAlsoLine = "### SEE ALSO"
	if bytes.Contains(p, []byte(seeAlsoLine)) {
		s.replace = true
		n := len(p)
		headerIdx := bytes.Index(p, []byte(seeAlsoLine))
		p, temp := p[:headerIdx], p[headerIdx:]
		for _, subCmd := range subcommands.All() {
			temp = bytes.ReplaceAll(temp, []byte(cmd.Name+" "+subCmd.Name()), []byte(subCmd.Name()))
			temp = bytes.ReplaceAll(temp, []byte(cmd.Name+"_"+subCmd.Name()), []byte(subCmd.Name()))
		}
		_, err := s.w.Write(append(p, temp...))
		return n, err
	}
	return s.w.Write(p)
}
