//go:build unix && !(linux || darwin || dragonfly || freebsd || netbsd || openbsd)

package loadavg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const Supported = true

func (l *LoadAvg) Update() error {
	cmd := exec.Command("sysctl", "-n", "vm.loadavg")
	cmd.Stderr = os.Stderr
	data, err := cmd.Output()
	if err != nil {
		return err
	}

	data = bytes.TrimPrefix(data, []byte("{ "))
	data = bytes.TrimSuffix(data, []byte(" }"))

	parts := strings.Fields(string(data))
	if len(parts) < 3 {
		return fmt.Errorf("%w: not enough values", ErrUnexpectedContent)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	for i, load := range parts[0:3] {
		if l.parts[i], err = strconv.ParseFloat(load, 64); err != nil {
			return err
		}
	}
	return err
}
