package loadavg

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	Supported = true
	path      = string(filepath.Separator) + "proc" + string(filepath.Separator) + "loadavg"
)

func (l *LoadAvg) Update() error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	l.mu.Lock()
	defer l.mu.Unlock()

	_, err = fmt.Fscanf(f, "%f %f %f", &l.min1, &l.min5, &l.min15)
	return err
}
