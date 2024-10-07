package ifdata

import (
	"errors"
	"fmt"
	"html/template"
	"net"
	"strings"

	"github.com/gabe565/moreutils/internal/cmdutil"
	"github.com/spf13/cobra"
)

const (
	Name                 = "ifdata"
	FlagExists           = "e"
	FlagPrint            = "p"
	FlagPrintExists      = "pe"
	FlagAddress          = "pa"
	FlagNetmask          = "pn"
	FlagNetworkAddress   = "pN"
	FlagBroadcastAddress = "pb"
	FlagMTU              = "pm"
	FlagFlags            = "pf"
	FlagHardwareAddress  = "ph"

	FlagInputStatistics  = "si"
	FlagInputPackets     = "sip"
	FlagInputBytes       = "sib"
	FlagInputErrors      = "sie"
	FlagInputDropped     = "sid"
	FlagInputFIFO        = "sif"
	FlagInputCompressed  = "sic"
	FlagInputMulticast   = "sim"
	FlagInputBytesSecond = "bips"

	FlagOutputStatistics    = "so"
	FlagOutputPackets       = "sop"
	FlagOutputBytes         = "sob"
	FlagOutputErrors        = "soe"
	FlagOutputDropped       = "sod"
	FlagOutputFIFO          = "sof"
	FlagOutputCollisions    = "sox"
	FlagOutputCarrierLosses = "soc"
	FlagOutputMulticast     = "som"
	FlagOutputBytesSecond   = "bops"
)

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " [flags] interface",
		Short:   "Get network interface info without parsing ifconfig output",
		RunE:    run,
		GroupID: cmdutil.Applet,
	}

	cmd.SetUsageFunc(usageFunc)
	cmd.DisableFlagsInUseLine = true
	cmd.DisableFlagParsing = true

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

var (
	ErrNoFormatter           = errors.New("no formatter was provided")
	ErrUnknownFormatter      = errors.New("unknown formatter")
	ErrNoInterface           = errors.New("no interface was provided")
	ErrInterfaceMissing      = errors.New("interface missing from /proc/net/dev")
	ErrStatisticsUnsupported = errors.New("platform does not support interface statistics")
)

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Usage()
	}

	var op, name string
	for _, arg := range args {
		switch {
		case arg == "-h", arg == "--help":
			return cmd.Usage()
		case arg == "-v", arg == "--version":
			if cmd.Version != "" {
				tmpl, err := template.New("").Parse(cmd.VersionTemplate())
				if err != nil {
					panic(err)
				}

				return tmpl.Execute(cmd.OutOrStdout(), cmd)
			}
		case strings.HasPrefix(arg, "-"):
			op = strings.TrimPrefix(arg, "-")
		default:
			name = arg
		}
	}

	switch {
	case op == "":
		return ErrNoFormatter
	case name == "":
		return ErrNoInterface
	}
	cmd.SilenceUsage = true

	iface, err := net.InterfaceByName(name)

	if op == FlagPrintExists {
		if err != nil {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "no")
		} else {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "yes")
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("%w: %s", err, name)
	}

	if op == FlagExists {
		return nil
	}

	switch op {
	case FlagMTU:
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), iface.MTU)
	case FlagFlags:
		for _, flag := range strings.Split(iface.Flags.String(), "|") {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), flag)
		}
	case FlagHardwareAddress:
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), strings.ToUpper(iface.HardwareAddr.String()))
	case FlagAddress, FlagNetmask, FlagNetworkAddress, FlagBroadcastAddress, FlagPrint:
		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			if addr, ok := addr.(*net.IPNet); ok && addr.IP.To4() != nil {
				switch op {
				case FlagAddress:
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), addr.IP)
				case FlagNetmask:
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), net.IP(addr.Mask))
				case FlagNetworkAddress:
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), addr.IP.Mask(addr.Mask))
				case FlagBroadcastAddress:
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), getBroadcastAddr(addr))
				case FlagPrint:
					_, _ = fmt.Fprintf(cmd.OutOrStdout(),
						"%s %s %s %d\n",
						addr.IP,
						net.IP(addr.Mask).String(),
						getBroadcastAddr(addr),
						iface.MTU,
					)
				}
			}
		}
	default:
		if statisticsSupported {
			return statistics(cmd, op, iface)
		}

		cmd.SilenceUsage = false
		return fmt.Errorf("%w: -%s", ErrUnknownFormatter, op)
	}

	return nil
}

func getBroadcastAddr(addr *net.IPNet) net.IP {
	if ip := addr.IP.To4(); ip != nil {
		mask := net.IP(addr.Mask)
		for i, b := range mask {
			ip[i] |= ^b
		}
		return ip
	}
	return nil
}
