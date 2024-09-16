package combine

import (
	"bufio"
	"io"
)

// iterLines reads lines from an io.Reader and calls the provided function.
func iterLines(r io.Reader, fn func(string) error) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := fn(scanner.Text()); err != nil {
			return err
		}
	}
	return scanner.Err()
}

// collectLines returns a map of lines from an io.Reader
func collectLines(r io.Reader) (map[string]struct{}, error) {
	seen := make(map[string]struct{})
	err := iterLines(r, func(line string) error {
		seen[line] = struct{}{}
		return nil
	})
	return seen, err
}
