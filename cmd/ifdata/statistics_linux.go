package ifdata

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/prometheus/procfs"
	"github.com/spf13/cobra"
)

const statisticsSupported = true

func statistics(cmd *cobra.Command, op string, iface *net.Interface) error {
	handle := statsHandler(op)
	if handle == nil {
		cmd.SilenceUsage = false
		return fmt.Errorf("%w: -%s", ErrUnknownFormatter, op)
	}

	device, err := getNetDevLine(iface.Name)
	if err != nil {
		return err
	}

	_, err = handle(cmd.OutOrStdout(), device)
	return err
}

func getNetDevLine(name string) (*procfs.NetDevLine, error) {
	fs, err := procfs.NewDefaultFS()
	if err != nil {
		return nil, err
	}

	entries, err := fs.NetDev()
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.Name == name {
			return &entry, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrInterfaceMissing, name)
}

func statsHandler(op string) func(w io.Writer, d *procfs.NetDevLine) (int, error) {
	switch op {
	case FlagInputStatistics:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) {
			return fmt.Fprintf(w, "%d %d %d %d %d %d %d %d\n",
				d.RxBytes, d.RxPackets,
				d.RxErrors, d.RxDropped,
				d.RxFIFO, d.RxFrame,
				d.RxCompressed, d.RxMulticast,
			)
		}
	case FlagInputPackets:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.RxPackets) }
	case FlagInputBytes:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.RxBytes) }
	case FlagInputErrors:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.RxErrors) }
	case FlagInputDropped:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.RxDropped) }
	case FlagInputFIFO:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.RxFIFO) }
	case FlagInputCompressed:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.RxCompressed) }
	case FlagInputMulticast:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.RxMulticast) }
	case FlagInputBytesSecond:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) {
			time.Sleep(time.Second)

			d2, err := getNetDevLine(d.Name)
			if err != nil {
				return 0, err
			}

			return fmt.Fprintln(w, d2.RxBytes-d.RxBytes)
		}

	case FlagOutputStatistics:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) {
			return fmt.Fprintf(w, "%d %d %d %d %d %d %d %d\n",
				d.TxBytes, d.TxPackets,
				d.TxErrors, d.TxDropped,
				d.TxFIFO, d.TxCollisions,
				d.TxCarrier, 0,
			)
		}
	case FlagOutputPackets:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.TxPackets) }
	case FlagOutputBytes:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.TxBytes) }
	case FlagOutputErrors:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.TxErrors) }
	case FlagOutputDropped:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.TxDropped) }
	case FlagOutputFIFO:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.TxFIFO) }
	case FlagOutputCollisions:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.TxCollisions) }
	case FlagOutputCarrierLosses:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.TxCarrier) }
	case FlagOutputMulticast:
		return func(w io.Writer, _ *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, 0) }
	case FlagOutputBytesSecond:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) {
			time.Sleep(time.Second)

			d2, err := getNetDevLine(d.Name)
			if err != nil {
				return 0, err
			}

			return fmt.Fprintln(w, d2.TxBytes-d.TxBytes)
		}

	default:
		return nil
	}
}
