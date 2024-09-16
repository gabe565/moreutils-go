package util

import (
	"os"

	"github.com/mattn/go-isatty"
)

// IsTerminal checks if an io.Reader or io.Writer is a TTY.
func IsTerminal(v any) bool {
	if f, ok := v.(*os.File); ok {
		return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
	}
	return false
}
