package ifdata

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/prometheus/procfs"
)

const statisticsSupported = true

func (f formatter) formatStatistics(iface *net.Interface) (string, error) {
	format := statsFormatter(f)
	if format == nil {
		return "", fmt.Errorf("%w: %s", ErrUnknownFormatter, f)
	}

	device, err := getNetDevLine(iface.Name)
	if err != nil {
		return "", err
	}

	return format(device)
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

func statsFormatter(op formatter) func(d *procfs.NetDevLine) (string, error) {
	switch op {
	case fmtInputStatistics:
		return func(d *procfs.NetDevLine) (string, error) {
			return fmt.Sprintf("%d %d %d %d %d %d %d %d",
				d.RxBytes, d.RxPackets,
				d.RxErrors, d.RxDropped,
				d.RxFIFO, d.RxFrame,
				d.RxCompressed, d.RxMulticast,
			), nil
		}
	case fmtInputPackets:
		return func(d *procfs.NetDevLine) (string, error) { return strconv.FormatUint(d.RxPackets, 10), nil }
	case fmtInputBytes:
		return func(d *procfs.NetDevLine) (string, error) { return strconv.FormatUint(d.RxBytes, 10), nil }
	case fmtInputErrors:
		return func(d *procfs.NetDevLine) (string, error) { return strconv.FormatUint(d.RxErrors, 10), nil }
	case fmtInputDropped:
		return func(d *procfs.NetDevLine) (string, error) { return strconv.FormatUint(d.RxDropped, 10), nil }
	case fmtInputFIFO:
		return func(d *procfs.NetDevLine) (string, error) { return strconv.FormatUint(d.RxFIFO, 10), nil }
	case fmtInputCompressed:
		return func(d *procfs.NetDevLine) (string, error) { return strconv.FormatUint(d.RxCompressed, 10), nil }
	case fmtInputMulticast:
		return func(d *procfs.NetDevLine) (string, error) { return strconv.FormatUint(d.RxMulticast, 10), nil }
	case fmtInputBytesSecond:
		return func(d *procfs.NetDevLine) (string, error) {
			time.Sleep(time.Second)

			d2, err := getNetDevLine(d.Name)
			if err != nil {
				return "", err
			}

			return strconv.FormatUint(d2.RxBytes-d.RxBytes, 10), nil
		}

	case fmtOutputStatistics:
		return func(d *procfs.NetDevLine) (string, error) {
			return fmt.Sprintf("%d %d %d %d %d %d %d %d",
				d.TxBytes, d.TxPackets,
				d.TxErrors, d.TxDropped,
				d.TxFIFO, d.TxCollisions,
				d.TxCarrier, 0,
			), nil
		}
	case fmtOutputPackets:
		return func(d *procfs.NetDevLine) (string, error) { return strconv.FormatUint(d.TxPackets, 10), nil }
	case fmtOutputBytes:
		return func(d *procfs.NetDevLine) (string, error) { return strconv.FormatUint(d.TxBytes, 10), nil }
	case fmtOutputErrors:
		return func(d *procfs.NetDevLine) (string, error) { return strconv.FormatUint(d.TxErrors, 10), nil }
	case fmtOutputDropped:
		return func(d *procfs.NetDevLine) (string, error) { return strconv.FormatUint(d.TxDropped, 10), nil }
	case fmtOutputFIFO:
		return func(d *procfs.NetDevLine) (string, error) { return strconv.FormatUint(d.TxFIFO, 10), nil }
	case fmtOutputCollisions:
		return func(d *procfs.NetDevLine) (string, error) { return strconv.FormatUint(d.TxCollisions, 10), nil }
	case fmtOutputCarrierLosses:
		return func(d *procfs.NetDevLine) (string, error) { return strconv.FormatUint(d.TxCarrier, 10), nil }
	case fmtOutputMulticast:
		return func(_ *procfs.NetDevLine) (string, error) { return "0", nil }
	case fmtOutputBytesSecond:
		return func(d *procfs.NetDevLine) (string, error) {
			time.Sleep(time.Second)

			d2, err := getNetDevLine(d.Name)
			if err != nil {
				return "", err
			}

			return strconv.FormatUint(d2.TxBytes-d.TxBytes, 10), nil
		}

	default:
		return nil
	}
}
