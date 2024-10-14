package editor

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/google/shlex"
	"github.com/mattn/go-tty"
)

const (
	envEditor     = "EDITOR"
	envVisual     = "VISUAL"
	defaultEditor = "vim"
)

var ErrUnset = errors.New("env is not set")

// Get checks VISUAL and EDITOR and returns the first result.
// If neither is set, it will return "vim"
func Get() ([]string, error) {
	var errs []error
	for _, env := range []string{envEditor, envVisual} {
		if editor, err := parseEnv(env); err == nil {
			return editor, errors.Join(errs...)
		} else if !errors.Is(err, ErrUnset) {
			errs = append(errs, err)
		}
	}

	return []string{defaultEditor}, errors.Join(errs...)
}

func parseEnv(env string) ([]string, error) {
	if val := os.Getenv(env); val != "" {
		editor, err := shlex.Split(val)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", env, err)
		}
		if len(editor) != 0 {
			return editor, nil
		}
	}
	return nil, fmt.Errorf("parse %s: %w", env, ErrUnset)
}

// Edit opens the configured editor with the given path.
// If forceTTY is true, "/dev/tty" will be opened for stdin and stdout.
func Edit(path string, forceTTY bool) error {
	editor, err := Get()
	if err != nil {
		slog.Warn("Failed to parse editor envs", "error", err)
	}

	editor = append(editor, path)

	cmd := exec.Command(editor[0], editor[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if forceTTY {
		t, err := tty.Open()
		if err != nil {
			return err
		}
		defer func() {
			_ = t.Close()
		}()

		cmd.Stdin = t.Input()
		cmd.Stdout = t.Output()
	}
	return cmd.Run()
}
