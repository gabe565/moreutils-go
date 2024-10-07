package ifdata

import (
	"strings"

	"github.com/spf13/cobra"
)

const Usage = `
  -h     help for ` + Name + `
  -e     Test to see if the interface exists, exit nonzero if it does not
  -p     Prints out the whole configuration of the interface
  -pe    Prints "yes" or "no" if the interface exists or not.
  -pa    Prints the IP address of the interface
  -pn    Prints the netmask of the interface
  -pN    Prints the network address of the interface
  -pb    Prints the broadcast address of the interface
  -pm    Prints the MTU of the interface
  -ph    Prints the hardware address of the interface. Exit with a failure exit code if there is not hardware address for the given network interface
  -pf    Prints the flags of the interface`

const UsageStatistics = `
  -si    Prints all input statistics of the interface
  -sip   Prints the number of input packets
  -sib   Prints the number of input bytes
  -sie   Prints the number of input errors
  -sid   Prints the number of dropped input packets
  -sif   Prints the number of input fifo overruns
  -sic   Prints the number of compressed input packets
  -sim   Prints the number of input multicast packets
  -so    Prints all output statistics of the interface
  -sop   Prints the number of output packets
  -sob   Prints the number of output bytes
  -soe   Prints the number of output errors
  -sod   Prints the number of dropped output packets
  -sof   Prints the number of output fifo overruns
  -sox   Prints the number of output collisions
  -soc   Prints the number of output carrier losses
  -som   Prints the number of output multicast packets
  -bips  Prints the number of bytes of incoming traffic measured in one second
  -bops  Prints the number of bytes of outgoing traffic measured in one second`

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

	newUsage := Usage
	if statisticsSupported {
		newUsage += UsageStatistics
	}
	if f := cmd.Flags().Lookup("version"); f != nil {
		newUsage += "\n  -" + f.Shorthand + "     " + f.Usage
	}

	tmpl = tmpl[:flagsIdx] + newUsage + tmpl[endIdx:]
	cmd.SetUsageTemplate(tmpl)
	cmd.SetUsageFunc(nil)
	return cmd.Usage()
}
