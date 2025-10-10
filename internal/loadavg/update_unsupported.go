//go:build !unix

package loadavg

import (
	"context"
	"fmt"
	"runtime"
)

const Supported = false

func (l *LoadAvg) Update(_ context.Context) error {
	return fmt.Errorf("%w: %s", ErrUnsupported, runtime.GOOS)
}
