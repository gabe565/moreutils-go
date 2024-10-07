package ifdata

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func UsageString(cmd *cobra.Command, full bool) string {
	var s strings.Builder
	s.WriteString("\n  -h      help for " + Name)

	for _, val := range formatterValues() {
		if val == fmtNone {
			continue
		}

		if full || val.supported() {
			_, _ = fmt.Fprintf(&s, "\n  %-7s %s", val, val.description())
		}
	}

	if cmd != nil {
		if f := cmd.Flags().Lookup("version"); f != nil {
			_, _ = fmt.Fprintf(&s, "\n  %-7s %s", "-"+f.Shorthand, f.Usage)
		}
	}

	return s.String()
}

func usageFunc(cmd *cobra.Command) error {
	tmpl := cmd.UsageTemplate()

	const flags = "Flags:"
	flagsIdx := strings.Index(tmpl, flags)
	if flagsIdx == -1 {
		panic("Missing usage flag start index")
	}
	flagsIdx += len(flags)

	const end = "{{end}}"
	endIdx := strings.Index(tmpl[flagsIdx:], end)
	if endIdx == -1 {
		panic("Missing usage flag end index")
	}
	endIdx += flagsIdx

	tmpl = tmpl[:flagsIdx] + UsageString(cmd, false) + tmpl[endIdx:]
	cmd.SetUsageTemplate(tmpl)
	cmd.SetUsageFunc(nil)
	return cmd.Usage()
}

func ManOptions() string {
	var s strings.Builder
	for _, val := range formatterValues() {
		if val == fmtNone {
			continue
		}

		_, _ = fmt.Fprintf(&s, "\n.PP\n\\fB%s\\fP\n\t%s\n", val, val.description())
	}
	return s.String()
}
