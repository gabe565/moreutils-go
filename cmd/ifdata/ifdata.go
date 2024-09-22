package ifdata

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/gabe565/moreutils/internal/cmdutil"
	"github.com/gabe565/moreutils/internal/util"
	"github.com/spf13/cobra"
)

const (
	Name                 = "ifdata"
	FlagPrint            = "print"
	FlagExists           = "exists"
	FlagAddress          = "address"
	FlagNetmask          = "netmask"
	FlagNetworkAddress   = "network-address"
	FlagBroadcastAddress = "broadcast-address"
	FlagMTU              = "mtu"
	FlagFlags            = "flags"
	FlagHardwareAddress  = "hardware-addr"
)

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " interface",
		Short:   "Get network interface info without parsing ifconfig output",
		RunE:    run,
		GroupID: cmdutil.Applet,
	}

	cmd.SetUsageFunc(usageFunc)

	cmd.Flags().Bool("help", false, "Print usage")
	cmd.Flags().Lookup("help").Hidden = true

	cmd.Flags().BoolP(FlagPrint, "p", false, "Prints out the whole configuration of the interface")
	cmd.Flags().BoolP(FlagExists, "e", false, "Test to see if the interface exists, exit nonzero if it does not")
	cmd.Flags().BoolP(FlagAddress, "a", false, "Prints the IPv4 address of the interface")
	cmd.Flags().BoolP(FlagNetmask, "n", false, "Prints the netmask of the interface")
	cmd.Flags().BoolP(FlagNetworkAddress, "N", false, "Prints the network address of the interface")
	cmd.Flags().BoolP(FlagBroadcastAddress, "b", false, "Prints the broadcast address of the interface")
	cmd.Flags().BoolP(FlagMTU, "m", false, "Prints the MTU of the interface")
	cmd.Flags().BoolP(FlagFlags, "f", false, "Prints the flags of the interface")
	cmd.Flags().BoolP(FlagHardwareAddress, "h", false, "Prints the hardware address of the interface. Exit with a failure exit code if there is not hardware address for the given network interface")

	cmd.MarkFlagsMutuallyExclusive(
		FlagExists,
		FlagAddress,
		FlagNetmask,
		FlagNetworkAddress,
		FlagBroadcastAddress,
		FlagMTU,
		FlagFlags,
		FlagHardwareAddress,
	)

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

var ErrNoOperation = errors.New("no operation was provided")

func run(cmd *cobra.Command, args []string) error {
	printFlag := util.Must2(cmd.Flags().GetBool(FlagPrint))
	exists := util.Must2(cmd.Flags().GetBool(FlagExists))
	address := util.Must2(cmd.Flags().GetBool(FlagAddress))
	netmask := util.Must2(cmd.Flags().GetBool(FlagNetmask))
	networkAddress := util.Must2(cmd.Flags().GetBool(FlagNetworkAddress))
	broadcastAddress := util.Must2(cmd.Flags().GetBool(FlagBroadcastAddress))
	mtu := util.Must2(cmd.Flags().GetBool(FlagMTU))
	flags := util.Must2(cmd.Flags().GetBool(FlagFlags))
	hardwareAddress := util.Must2(cmd.Flags().GetBool(FlagHardwareAddress))

	if len(args) == 0 {
		// Error is suppressed when "-h" is provided with no flags
		if hardwareAddress || cmd.Flags().NFlag() == 0 {
			return cmd.Usage()
		}
		return cobra.ExactArgs(1)(cmd, args)
	}
	cmd.SilenceUsage = true

	iface, err := net.InterfaceByName(args[0])

	if printFlag && exists {
		if err != nil {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "no")
		} else {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "yes")
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("%w: %s", err, args[0])
	}

	if exists {
		return nil
	}

	switch {
	case mtu:
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), iface.MTU)
	case flags:
		for _, flag := range strings.Split(iface.Flags.String(), "|") {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), flag)
		}
	case hardwareAddress:
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), strings.ToUpper(iface.HardwareAddr.String()))
	case address, netmask, networkAddress, broadcastAddress, printFlag:
		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			if addr, ok := addr.(*net.IPNet); ok && addr.IP.To4() != nil {
				switch {
				case address:
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), addr.IP)
				case netmask:
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), net.IP(addr.Mask))
				case networkAddress:
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), addr.IP.Mask(addr.Mask))
				case broadcastAddress:
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), getBroadcastAddr(addr))
				case printFlag:
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
		cmd.SilenceUsage = false
		return ErrNoOperation
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
