package seealsoreplacer

import (
	"bytes"
	"io"
	"regexp"
	"slices"

	"gabe565.com/moreutils/cmd"
	"github.com/spf13/cobra"
)

func New(w io.Writer, search string, cmds []*cobra.Command) *SeeAlsoReplacer {
	return &SeeAlsoReplacer{
		w:      w,
		search: []byte(search),
		cmds:   cmds,
	}
}

type SeeAlsoReplacer struct {
	w      io.Writer
	search []byte
	cmds   []*cobra.Command
}

func (s *SeeAlsoReplacer) Write(p []byte) (int, error) {
	if bytes.Contains(p, s.search) {
		i := bytes.Index(p, s.search)
		if i == -1 {
			panic("Missing header index")
		}
		i += len(s.search)

		section := slices.Clone(p[i:])
		for _, subCmd := range s.cmds {
			re := regexp.MustCompile(regexp.QuoteMeta(cmd.Name) + `[ _-]` + regexp.QuoteMeta(subCmd.Name()))
			section = re.ReplaceAll(section, []byte(subCmd.Name()))
		}

		_, err := s.w.Write(slices.Concat(p[:i], section))
		return len(p), err
	}
	return s.w.Write(p)
}
