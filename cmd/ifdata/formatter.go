package ifdata

import (
	"errors"
	"fmt"
	"net"
	"slices"
	"strconv"
	"strings"
)

//go:generate go run github.com/dmarkham/enumer -type formatter -linecomment -output formatter_string.go

type formatter uint8

const (
	fmtNone formatter = iota //

	fmtExists           // -e
	fmtPrint            // -p
	fmtPrintExists      // -pe
	fmtAddress          // -pa
	fmtNetworkAddress   // -pN
	fmtNetmask          // -pn
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
	case fmtNetworkAddress:
		return "Prints the network address of the interface"
	case fmtNetmask:
		return "Prints the netmask of the interface"
	case fmtBroadcastAddress:
		return "Prints the broadcast address of the interface"
	case fmtMTU:
		return "Prints the MTU of the interface"
	case fmtFlags:
		return "Prints the flags of the interface"
	case fmtHardwareAddress:
		return "Prints the hardware address of the interface. Exit with a failure exit code if there is not hardware address for the given network interface"
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
	case fmtOutputStatistics:
		return "Prints all output statistics of the interface"
	case fmtOutputPackets:
		return "Prints the number of output packets"
	case fmtOutputBytes:
		return "Prints the number of output bytes"
	case fmtOutputErrors:
		return "Prints the number of output errors"
	case fmtOutputDropped:
		return "Prints the number of dropped output packets"
	case fmtOutputFIFO:
		return "Prints the number of output fifo overruns"
	case fmtOutputCollisions:
		return "Prints the number of output collisions"
	case fmtOutputCarrierLosses:
		return "Prints the number of output carrier losses"
	case fmtOutputMulticast:
		return "Prints the number of output multicast packets"
	case fmtInputBytesSecond:
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

var ErrUnknownFormatter = errors.New("unknown formatter")

func (f formatter) Sprint(iface *net.Interface) (string, error) {
	switch f {
	case fmtMTU:
		return strconv.Itoa(iface.MTU), nil
	case fmtFlags:
		return strings.Join(strings.Split(iface.Flags.String(), "|"), "\n"), nil
	case fmtHardwareAddress:
		return strings.ToUpper(iface.HardwareAddr.String()), nil
	case fmtAddress, fmtNetmask, fmtNetworkAddress, fmtBroadcastAddress, fmtPrint:
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		var buf strings.Builder
		for _, addr := range addrs {
			if addr, ok := addr.(*net.IPNet); ok && addr.IP.To4() != nil {
				switch f {
				case fmtAddress:
					_, _ = fmt.Fprintln(&buf, addr.IP)
				case fmtNetmask:
					_, _ = fmt.Fprintln(&buf, net.IP(addr.Mask))
				case fmtNetworkAddress:
					_, _ = fmt.Fprintln(&buf, addr.IP.Mask(addr.Mask))
				case fmtBroadcastAddress:
					_, _ = fmt.Fprintln(&buf, getBroadcastAddr(addr))
				case fmtPrint:
					_, _ = fmt.Fprintf(&buf,
						"%s %s %s %d\n",
						addr.IP,
						net.IP(addr.Mask).String(),
						getBroadcastAddr(addr),
						iface.MTU,
					)
				}
			}
		}
		return strings.TrimSpace(buf.String()), nil
	default:
		if f.supported() {
			return f.formatStatistics(iface)
		}
		return "", fmt.Errorf("%w: %s", ErrUnknownFormatter, f)
	}
}

func getBroadcastAddr(addr *net.IPNet) net.IP {
	if ip := addr.IP.To4(); ip != nil {
		ip := slices.Clone(ip)
		for i, b := range addr.Mask {
			ip[i] |= ^b
		}
		return ip
	}
	return nil
}
