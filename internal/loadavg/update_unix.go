//go:build unix && !(linux || darwin || dragonfly || freebsd || netbsd || openbsd)

package loadavg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

const Supported = true

func (l *LoadAvg) Update() error {
	cmd := exec.Command("sysctl", "-n", "vm.loadavg")
	cmd.Stderr = os.Stderr
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	_, err = fmt.Fscanf(bytes.NewReader(stdout), "{ %f %f %f }", &l.min1, &l.min5, &l.min15)
	return err
}
