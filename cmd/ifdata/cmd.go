package ifdata

import (
	"errors"
	"fmt"
	"html/template"
	"io"
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
	ErrNoInterface           = errors.New("no interface was provided")
	ErrInterfaceMissing      = errors.New("interface missing from /proc/net/dev")
	ErrStatisticsUnsupported = errors.New("platform does not support interface statistics")
)

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	var format formatter
	var name string
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
			name = arg
		}
	}

	switch {
	case format == fmtNone:
		return ErrNoFormatter
	case name == "":
		return ErrNoInterface
	}
	cmd.SilenceUsage = true

	iface, err := net.InterfaceByName(name)

	if format == fmtPrintExists {
		if err != nil {
			_, _ = io.WriteString(cmd.OutOrStdout(), "no\n")
		} else {
			_, _ = io.WriteString(cmd.OutOrStdout(), "yes\n")
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("%w: %s", err, name)
	}

	if format == fmtExists {
		return nil
	}

	s, err := format.Sprint(cmd, iface)
	if err != nil {
		return err
	}

	_, err = io.WriteString(cmd.OutOrStdout(), s+"\n")
	return err
}
