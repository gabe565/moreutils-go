//go:build darwin || dragonfly || freebsd || netbsd || openbsd

package loadavg

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

const Supported = true

// From https://github.com/prometheus/node_exporter/blob/master/collector/loadavg_bsd.go

func (l *LoadAvg) Update() error {
	type loadavg struct {
		load  [3]uint32
		scale int
	}

	b, err := unix.SysctlRaw("vm.loadavg")
	if err != nil {
		return err
	}

	load := *(*loadavg)(unsafe.Pointer(&b[0]))
	scale := float64(load.scale)

	l.mu.Lock()
	defer l.mu.Unlock()

	l.min1 = float64(load.load[0]) / scale
	l.min5 = float64(load.load[1]) / scale
	l.min15 = float64(load.load[2]) / scale
	return nil
}
