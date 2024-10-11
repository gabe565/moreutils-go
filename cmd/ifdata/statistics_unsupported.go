//go:build !(linux || darwin)

package ifdata

import (
	"net"

	"github.com/spf13/cobra"
)

const statisticsSupported = false

func (f formatter) formatStatistics(_ *cobra.Command, _ *net.Interface) (string, error) {
	return "", ErrStatisticsUnsupported
}
