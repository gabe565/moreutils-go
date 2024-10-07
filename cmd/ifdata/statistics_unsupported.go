//go:build !(linux || darwin)

package ifdata

import (
	"net"

	"github.com/spf13/cobra"
)

const statisticsSupported = false

func statistics(_ *cobra.Command, _ string, _ *net.Interface) error {
	return ErrStatisticsUnsupported
}
