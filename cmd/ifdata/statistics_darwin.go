package ifdata

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"time"

	"golang.org/x/sys/unix"
)

const statisticsSupported = true

func (f formatter) formatStatistics(iface *net.Interface) (string, error) {
	format := statsFormatter(f)
	if format == nil {
		return "", fmt.Errorf("%w: %s", ErrUnknownFormatter, f)
	}

	data, err := getIfaceData(iface.Index)
	if err != nil {
		return "", err
	}

	return format(data)
}

func statsFormatter(op formatter) func(d *ifMsghdr2) (string, error) {
	switch op {
	case fmtInputStatistics:
		return func(d *ifMsghdr2) (string, error) {
			return fmt.Sprintf("%d %d %d %d %d %d %d %d",
				d.Data.Ibytes, d.Data.Ipackets,
				d.Data.Ierrors, d.Data.Iqdrops,
				0, 0,
				0, d.Data.Imcasts,
			), nil
		}
	case fmtInputPackets:
		return func(d *ifMsghdr2) (string, error) { return strconv.FormatUint(d.Data.Ipackets, 10), nil }
	case fmtInputBytes:
		return func(d *ifMsghdr2) (string, error) { return strconv.FormatUint(d.Data.Ibytes, 10), nil }
	case fmtInputErrors:
		return func(d *ifMsghdr2) (string, error) { return strconv.FormatUint(d.Data.Ierrors, 10), nil }
	case fmtInputDropped:
		return func(d *ifMsghdr2) (string, error) { return strconv.FormatUint(d.Data.Iqdrops, 10), nil }
	case fmtInputMulticast:
		return func(d *ifMsghdr2) (string, error) { return strconv.FormatUint(d.Data.Imcasts, 10), nil }
	case fmtInputFIFO, fmtInputCompressed:
		return func(_ *ifMsghdr2) (string, error) { return "0", nil }
	case fmtInputBytesSecond:
		return func(d *ifMsghdr2) (string, error) {
			time.Sleep(time.Second)

			d2, err := getIfaceData(int(d.Index))
			if err != nil {
				return "", err
			}

			return strconv.FormatUint(d2.Data.Ibytes-d.Data.Ibytes, 10), nil
		}

	case fmtOutputStatistics:
		return func(d *ifMsghdr2) (string, error) {
			return fmt.Sprintf("%d %d %d %d %d %d %d %d",
				d.Data.Obytes, d.Data.Opackets,
				d.Data.Oerrors, 0,
				0, d.Data.Collisions,
				0, d.Data.Omcasts,
			), nil
		}
	case fmtOutputPackets:
		return func(d *ifMsghdr2) (string, error) { return strconv.FormatUint(d.Data.Opackets, 10), nil }
	case fmtOutputBytes:
		return func(d *ifMsghdr2) (string, error) { return strconv.FormatUint(d.Data.Obytes, 10), nil }
	case fmtOutputErrors:
		return func(d *ifMsghdr2) (string, error) { return strconv.FormatUint(d.Data.Oerrors, 10), nil }
	case fmtOutputCollisions:
		return func(d *ifMsghdr2) (string, error) { return strconv.FormatUint(d.Data.Collisions, 10), nil }
	case fmtOutputMulticast:
		return func(d *ifMsghdr2) (string, error) { return strconv.FormatUint(d.Data.Omcasts, 10), nil }
	case fmtOutputDropped, fmtOutputFIFO, fmtOutputCarrierLosses:
		return func(_ *ifMsghdr2) (string, error) { return "0", nil }
	case fmtOutputBytesSecond:
		return func(d *ifMsghdr2) (string, error) {
			time.Sleep(time.Second)

			d2, err := getIfaceData(int(d.Index))
			if err != nil {
				return "", err
			}

			return strconv.FormatUint(d2.Data.Obytes-d.Data.Obytes, 10), nil
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
