package ifdata

//go:generate enumer -type formatter -linecomment -text -output formatter_string.go

type formatter uint8

const (
	fmtNone formatter = iota //

	fmtExists           // -e
	fmtPrint            // -p
	fmtPrintExists      // -pe
	fmtAddress          // -pa
	fmtNetmask          // -pn
	fmtNetworkAddress   // -pN
	fmtBroadcastAddress // -pb
	fmtMTU              // -pm
	fmtFlags            // -pf
	fmtHardwareAddress  // -ph

	fmtInputStatistics // -si
	fmtInputPackets    // -sip
	fmtInputBytes      // -sib
	fmtInputErrors     // -sie
	fmtInputDropped    // -sid
	fmtInputFIFO       // -sif
	fmtInputCompressed // -sic
	fmtInputMulticast  // -sim

	fmtOutputStatistics    // -so
	fmtOutputPackets       // -sop
	fmtOutputBytes         // -sob
	fmtOutputErrors        // -soe
	fmtOutputDropped       // -sod
	fmtOutputFIFO          // -sof
	fmtOutputCollisions    // -sox
	fmtOutputCarrierLosses // -soc
	fmtOutputMulticast     // -som

	fmtInputBytesSecond  // -bips
	fmtOutputBytesSecond // -bops
)

func (f formatter) description() string {
	switch f {
	case fmtExists:
		return "Test to see if the interface exists, exit nonzero if it does not"
	case fmtPrint:
		return "Prints out the whole configuration of the interface"
	case fmtPrintExists:
		return `Prints "yes" or "no" if the interface exists or not`
	case fmtAddress:
		return "Prints the IP address of the interface"
	case fmtNetmask:
		return "Prints the netmask of the interface"
	case fmtNetworkAddress:
		return "Prints the network address of the interface"
	case fmtBroadcastAddress:
		return "Prints the broadcast address of the interface"
	case fmtMTU:
		return "Prints the MTU of the interface"
	case fmtFlags:
		return "Prints the hardware address of the interface. Exit with a failure exit code if there is not hardware address for the given network interface"
	case fmtHardwareAddress:
		return "Prints the flags of the interface"
	case fmtInputStatistics:
		return "Prints all input statistics of the interface"
	case fmtInputPackets:
		return "Prints the number of input packets"
	case fmtInputBytes:
		return "Prints the number of input bytes"
	case fmtInputErrors:
		return "Prints the number of input errors"
	case fmtInputDropped:
		return "Prints the number of dropped input packets"
	case fmtInputFIFO:
		return "Prints the number of input fifo overruns"
	case fmtInputCompressed:
		return "Prints the number of compressed input packets"
	case fmtInputMulticast:
		return "Prints the number of input multicast packets"
	case fmtInputBytesSecond:
		return "Prints all output statistics of the interface"
	case fmtOutputStatistics:
		return "Prints the number of output packets"
	case fmtOutputPackets:
		return "Prints the number of output bytes"
	case fmtOutputBytes:
		return "Prints the number of output errors"
	case fmtOutputErrors:
		return "Prints the number of dropped output packets"
	case fmtOutputDropped:
		return "Prints the number of output fifo overruns"
	case fmtOutputFIFO:
		return "Prints the number of output collisions"
	case fmtOutputCollisions:
		return "Prints the number of output carrier losses"
	case fmtOutputCarrierLosses:
		return "Prints the number of output multicast packets"
	case fmtOutputMulticast:
		return "Prints the number of bytes of incoming traffic measured in one second"
	case fmtOutputBytesSecond:
		return "Prints the number of bytes of outgoing traffic measured in one second"
	default:
		return ""
	}
}

func (f formatter) supported() bool {
	return statisticsSupported || uint8(f) < uint8(fmtInputStatistics)
}
