package ifdata

import (
	"cmp"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net"
	"slices"
	"strings"

	"github.com/gabe565/moreutils/internal/cmdutil"
	"github.com/gabe565/moreutils/internal/util"
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
	ErrNoInterface           = errors.New("no interface was provided")
	ErrInterfaceMissing      = errors.New("interface missing from /proc/net/dev")
	ErrStatisticsUnsupported = errors.New("platform does not support interface statistics")
)

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	var format formatter
	names := make([]string, 0, len(args)-1)
	for _, arg := range args {
		switch {
		case arg == "-h", arg == "--help":
			return cmd.Help()
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
			if format, err = formatterString(arg); err != nil {
				return err
			}
		default:
			names = append(names, arg)
		}
	}
	args = names

	switch {
	case format == fmtNone:
		return ErrNoFormatter
	case len(args) == 0 && (format == fmtExists || format == fmtPrintExists):
		return ErrNoInterface
	}
	cmd.SilenceUsage = true

	ifaces := make([]*net.Interface, 0, len(args))
	var errs []error
	for _, arg := range args {
		iface, err := net.InterfaceByName(arg)
		if format == fmtPrintExists {
			var s string
			if len(args) != 1 {
				s += arg + " "
			}
			if err == nil {
				s += "yes\n"
			} else {
				s += "no\n"
			}
			_, _ = io.WriteString(cmd.OutOrStdout(), s)
		} else {
			if err == nil {
				ifaces = append(ifaces, iface)
			} else {
				errs = append(errs, fmt.Errorf("%w: %s", err, arg))
			}
		}
	}

	if format == fmtExists || format == fmtPrintExists {
		return errors.Join(errs...)
	}

	if len(ifaces) == 0 {
		v, err := net.Interfaces()
		if err != nil {
			return err
		}

		for _, iface := range v {
			ifaces = append(ifaces, &iface)
		}

		slices.SortStableFunc(ifaces, func(a, b *net.Interface) int {
			return cmp.Compare(a.Name, b.Name)
		})
	}

	for _, iface := range ifaces {
		s, err := format.Sprint(cmd, iface)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if s != "" {
			if len(args) != 1 {
				s = iface.Name + " " + s
			}
			_, _ = io.WriteString(cmd.OutOrStdout(), s+"\n")
		}
	}
	return util.JoinErrors(errs...)
}
