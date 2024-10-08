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

const Name = "ifdata"

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " [flags] interface",
		Short:   "Get network interface info without parsing ifconfig output",
		RunE:    run,
		GroupID: cmdutil.Applet,

		ValidArgsFunction: validArgs,
	}

	cmd.InitDefaultHelpFlag()

	cmd.SetUsageFunc(usageFunc)
	cmd.DisableFlagsInUseLine = true
	cmd.DisableFlagParsing = true

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func validArgs(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		vals := formatterValues()
		strs := make([]string, 0, len(vals))
		for _, val := range vals {
			if val == fmtNone {
				continue
			}

			if val.supported() {
				strs = append(strs, val.String()+"\t"+val.description())
			}
		}
		return strs, cobra.ShellCompDirectiveNoFileComp
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	names := make([]string, 0, len(ifaces))
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		desc := make([]string, 0, len(addrs))
		for _, addr := range addrs {
			if addr, ok := addr.(*net.IPNet); ok && addr.IP.To4() != nil {
				desc = append(desc, addr.String())
			}
		}
		names = append(names, iface.Name+"\t"+strings.Join(desc, ","))
	}
	return names, cobra.ShellCompDirectiveNoFileComp
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

	var op formatter
	var name string
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
			var err error
			if op, err = formatterString(arg); err != nil {
				return err
			}
		default:
			name = arg
		}
	}

	switch {
	case op == fmtNone:
		return ErrNoFormatter
	case name == "":
		return ErrNoInterface
	}
	cmd.SilenceUsage = true

	iface, err := net.InterfaceByName(name)

	if op == fmtPrintExists {
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

	if op == fmtExists {
		return nil
	}

	switch op {
	case fmtMTU:
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), iface.MTU)
	case fmtFlags:
		for _, flag := range strings.Split(iface.Flags.String(), "|") {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), flag)
		}
	case fmtHardwareAddress:
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), strings.ToUpper(iface.HardwareAddr.String()))
	case fmtAddress, fmtNetmask, fmtNetworkAddress, fmtBroadcastAddress, fmtPrint:
		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			if addr, ok := addr.(*net.IPNet); ok && addr.IP.To4() != nil {
				switch op {
				case fmtAddress:
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), addr.IP)
				case fmtNetmask:
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), net.IP(addr.Mask))
				case fmtNetworkAddress:
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), addr.IP.Mask(addr.Mask))
				case fmtBroadcastAddress:
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), getBroadcastAddr(addr))
				case fmtPrint:
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
		return fmt.Errorf("%w: %s", ErrUnknownFormatter, op)
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
