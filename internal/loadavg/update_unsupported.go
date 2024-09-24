//go:build !unix

package loadavg

import (
	"fmt"
	"runtime"
)

const Supported = false

func (l *LoadAvg) Update() error {
	return fmt.Errorf("%w: %s", ErrUnsupported, runtime.GOOS)
}
