package ifdata

import (
	"fmt"
	"io"
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
	cmd.SetUsageFunc(nil)
	u := cmd.UsageString()
	cmd.SetUsageFunc(usageFunc)

	const flags = "Flags:"
	flagsIdx := strings.Index(u, flags)
	if flagsIdx == -1 {
		panic("Missing usage flag start index")
	}
	flagsIdx += len(flags)

	u = u[:flagsIdx] + UsageString(cmd, false) + "\n"
	_, err := io.WriteString(cmd.OutOrStdout(), u)
	return err
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
