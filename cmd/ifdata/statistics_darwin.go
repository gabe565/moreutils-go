package ifdata

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

const statisticsSupported = true

func statistics(cmd *cobra.Command, op formatter, iface *net.Interface) error {
	handle := statsHandler(op)
	if handle == nil {
		cmd.SilenceUsage = false
		return fmt.Errorf("%w: %s", ErrUnknownFormatter, op)
	}

	data, err := getIfaceData(iface.Index)
	if err != nil {
		return err
	}

	_, err = handle(cmd.OutOrStdout(), data)
	return err
}

func statsHandler(op formatter) func(w io.Writer, d *ifMsghdr2) (int, error) {
	switch op {
	case fmtInputStatistics:
		return func(w io.Writer, d *ifMsghdr2) (int, error) {
			return fmt.Fprintf(w, "%d %d %d %d %d %d %d %d\n",
				d.Data.Ibytes, d.Data.Ipackets,
				d.Data.Ierrors, d.Data.Iqdrops,
				0, 0,
				0, d.Data.Imcasts,
			)
		}
	case fmtInputPackets:
		return func(w io.Writer, d *ifMsghdr2) (int, error) { return fmt.Fprintln(w, d.Data.Ipackets) }
	case fmtInputBytes:
		return func(w io.Writer, d *ifMsghdr2) (int, error) { return fmt.Fprintln(w, d.Data.Ibytes) }
	case fmtInputErrors:
		return func(w io.Writer, d *ifMsghdr2) (int, error) { return fmt.Fprintln(w, d.Data.Ierrors) }
	case fmtInputDropped:
		return func(w io.Writer, d *ifMsghdr2) (int, error) { return fmt.Fprintln(w, d.Data.Iqdrops) }
	case fmtInputMulticast:
		return func(w io.Writer, d *ifMsghdr2) (int, error) { return fmt.Fprintln(w, d.Data.Imcasts) }
	case fmtInputFIFO, fmtInputCompressed:
		return func(w io.Writer, _ *ifMsghdr2) (int, error) { return fmt.Fprintln(w, 0) }
	case fmtInputBytesSecond:
		return func(w io.Writer, d *ifMsghdr2) (int, error) {
			time.Sleep(time.Second)

			d2, err := getIfaceData(int(d.Index))
			if err != nil {
				return 0, err
			}

			return fmt.Fprintln(w, d2.Data.Ibytes-d.Data.Ibytes)
		}

	case fmtOutputStatistics:
		return func(w io.Writer, d *ifMsghdr2) (int, error) {
			return fmt.Fprintf(w, "%d %d %d %d %d %d %d %d\n",
				d.Data.Obytes, d.Data.Opackets,
				d.Data.Oerrors, 0,
				0, d.Data.Collisions,
				0, d.Data.Omcasts,
			)
		}
	case fmtOutputPackets:
		return func(w io.Writer, d *ifMsghdr2) (int, error) { return fmt.Fprintln(w, d.Data.Opackets) }
	case fmtOutputBytes:
		return func(w io.Writer, d *ifMsghdr2) (int, error) { return fmt.Fprintln(w, d.Data.Obytes) }
	case fmtOutputErrors:
		return func(w io.Writer, d *ifMsghdr2) (int, error) { return fmt.Fprintln(w, d.Data.Oerrors) }
	case fmtOutputCollisions:
		return func(w io.Writer, d *ifMsghdr2) (int, error) { return fmt.Fprintln(w, d.Data.Collisions) }
	case fmtOutputMulticast:
		return func(w io.Writer, d *ifMsghdr2) (int, error) { return fmt.Fprintln(w, d.Data.Omcasts) }
	case fmtOutputDropped, fmtOutputFIFO, fmtOutputCarrierLosses:
		return func(w io.Writer, _ *ifMsghdr2) (int, error) { return fmt.Fprintln(w, 0) }
	case fmtOutputBytesSecond:
		return func(w io.Writer, d *ifMsghdr2) (int, error) {
			time.Sleep(time.Second)

			d2, err := getIfaceData(int(d.Index))
			if err != nil {
				return 0, err
			}

			return fmt.Fprintln(w, d2.Data.Obytes-d.Data.Obytes)
		}

	default:
		return nil
	}
}

// From https://github.com/prometheus/node_exporter/blob/master/collector/netdev_darwin.go

func getIfaceData(index int) (*ifMsghdr2, error) {
	var data ifMsghdr2
	rawData, err := unix.SysctlRaw("net", unix.AF_ROUTE, 0, 0, unix.NET_RT_IFLIST2, index)
	if err != nil {
		return nil, err
	}

	err = binary.Read(bytes.NewReader(rawData), binary.LittleEndian, &data)
	return &data, err
}

type ifMsghdr2 struct {
	Msglen    uint16
	Version   uint8
	Type      uint8
	Addrs     int32
	Flags     int32
	Index     uint16
	_         [2]byte
	SndLen    int32
	SndMaxlen int32
	SndDrops  int32
	Timer     int32
	Data      ifData64
}

// https://github.com/apple/darwin-xnu/blob/main/bsd/net/if_var.h#L199-L231
type ifData64 struct {
	Type       uint8
	Typelen    uint8
	Physical   uint8
	Addrlen    uint8
	Hdrlen     uint8
	Recvquota  uint8
	Xmitquota  uint8
	Unused1    uint8
	Mtu        uint32
	Metric     uint32
	Baudrate   uint64
	Ipackets   uint64
	Ierrors    uint64
	Opackets   uint64
	Oerrors    uint64
	Collisions uint64
	Ibytes     uint64
	Obytes     uint64
	Imcasts    uint64
	Omcasts    uint64
	Iqdrops    uint64
	Noproto    uint64
	Recvtiming uint32
	Xmittiming uint32
	Lastchange unix.Timeval32
}
