//go:build !(linux || darwin)

package ifdata

import "net"

const statisticsSupported = false

func (f formatter) formatStatistics(_ *net.Interface) (string, error) {
	return "", ErrStatisticsUnsupported
}
