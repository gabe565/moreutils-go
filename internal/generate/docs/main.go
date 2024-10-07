package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gabe565/moreutils/cmd"
	"github.com/gabe565/moreutils/cmd/ifdata"
	"github.com/gabe565/moreutils/internal/cmdutil"
	"github.com/gabe565/moreutils/internal/cmdutil/subcommands"
	"github.com/gabe565/moreutils/internal/generate/seealsoreplacer"
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
			w = seealsoreplacer.New(w, "### SEE ALSO\n", subcommands.All())
		}

		if subCmd.Name() == ifdata.Name {
			w = optionsReplacer{
				w:       w,
				replace: ifdata.UsageString(subCmd, true),
			}
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

type optionsReplacer struct {
	w       io.Writer
	replace string
}

func (o optionsReplacer) Write(p []byte) (int, error) {
	if !strings.HasSuffix(o.replace, "\n") {
		o.replace += "\n"
	}

	const header, footer = "\n### Options\n\n```", "```"
	if bytes.Contains(p, []byte(header)) {
		beforeIdx := bytes.Index(p, []byte(header))
		if beforeIdx == -1 {
			panic("missing header: " + header)
		}

		afterIdx := bytes.Index(p[beforeIdx+len(header):], []byte(footer))
		if afterIdx == -1 {
			panic("missing footer: " + footer)
		}
		afterIdx += beforeIdx + len(header)

		_, err := o.w.Write(slices.Concat(p[:beforeIdx], []byte(header), []byte(o.replace), p[afterIdx:]))
		return len(p), err
	}
	return o.w.Write(p)
}
