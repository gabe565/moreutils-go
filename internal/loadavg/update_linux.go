package loadavg

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	Supported = true
	path      = string(filepath.Separator) + "proc" + string(filepath.Separator) + "loadavg"
)

func (l *LoadAvg) Update() error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	parts := strings.Fields(string(data))
	if len(parts) < 3 {
		return fmt.Errorf("%w in %s: not enough values", ErrUnexpectedContent, path)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	for i, load := range parts[0:3] {
		if l.parts[i], err = strconv.ParseFloat(load, 64); err != nil {
			return err
		}
	}
	return nil
}
