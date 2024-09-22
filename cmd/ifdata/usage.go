package ifdata

import (
	"strings"

	"github.com/spf13/cobra"
)

const Usage = `
  -h    help for ` + Name + `
  -e    Test to see if the interface exists, exit nonzero if it does not
  -p    Prints out the whole configuration of the interface
  -pe   Prints "yes" or "no" if the interface exists or not.
  -pa   Prints the IP address of the interface
  -pn   Prints the netmask of the interface
  -pN   Prints the network address of the interface
  -pb   Prints the broadcast address of the interface
  -pm   Prints the MTU of the interface
  -ph   Prints the hardware address of the interface. Exit with a failure exit code if there is not hardware address for the given network interface
  -pf   Prints the flags of the interface`

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

	tmpl = tmpl[:flagsIdx] + Usage + tmpl[endIdx:]
	cmd.SetUsageTemplate(tmpl)
	cmd.SetUsageFunc(nil)
	return cmd.Usage()
}
