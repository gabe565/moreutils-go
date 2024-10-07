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

func statistics(cmd *cobra.Command, op formatter, iface *net.Interface) error {
	handle := statsHandler(op)
	if handle == nil {
		cmd.SilenceUsage = false
		return fmt.Errorf("%w: %s", ErrUnknownFormatter, op)
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

func statsHandler(op formatter) func(w io.Writer, d *procfs.NetDevLine) (int, error) {
	switch op {
	case fmtInputStatistics:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) {
			return fmt.Fprintf(w, "%d %d %d %d %d %d %d %d\n",
				d.RxBytes, d.RxPackets,
				d.RxErrors, d.RxDropped,
				d.RxFIFO, d.RxFrame,
				d.RxCompressed, d.RxMulticast,
			)
		}
	case fmtInputPackets:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.RxPackets) }
	case fmtInputBytes:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.RxBytes) }
	case fmtInputErrors:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.RxErrors) }
	case fmtInputDropped:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.RxDropped) }
	case fmtInputFIFO:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.RxFIFO) }
	case fmtInputCompressed:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.RxCompressed) }
	case fmtInputMulticast:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.RxMulticast) }
	case fmtInputBytesSecond:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) {
			time.Sleep(time.Second)

			d2, err := getNetDevLine(d.Name)
			if err != nil {
				return 0, err
			}

			return fmt.Fprintln(w, d2.RxBytes-d.RxBytes)
		}

	case fmtOutputStatistics:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) {
			return fmt.Fprintf(w, "%d %d %d %d %d %d %d %d\n",
				d.TxBytes, d.TxPackets,
				d.TxErrors, d.TxDropped,
				d.TxFIFO, d.TxCollisions,
				d.TxCarrier, 0,
			)
		}
	case fmtOutputPackets:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.TxPackets) }
	case fmtOutputBytes:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.TxBytes) }
	case fmtOutputErrors:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.TxErrors) }
	case fmtOutputDropped:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.TxDropped) }
	case fmtOutputFIFO:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.TxFIFO) }
	case fmtOutputCollisions:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.TxCollisions) }
	case fmtOutputCarrierLosses:
		return func(w io.Writer, d *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, d.TxCarrier) }
	case fmtOutputMulticast:
		return func(w io.Writer, _ *procfs.NetDevLine) (int, error) { return fmt.Fprintln(w, 0) }
	case fmtOutputBytesSecond:
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
