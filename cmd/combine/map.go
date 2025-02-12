package combine

import (
	"bufio"
	"io"
	"iter"
)

// iterLines returns an iterator over each line of the provided io.Reader.
func iterLines(r io.Reader) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			if !yield(scanner.Text(), nil) {
				return
			}
		}
		if scanner.Err() != nil {
			yield("", scanner.Err())
		}
	}
}

// collectLines returns a map of lines from an io.Reader
func collectLines(r io.Reader) (map[string]struct{}, error) {
	seen := make(map[string]struct{})
	for line, err := range iterLines(r) {
		if err != nil {
			return nil, err
		}
		seen[line] = struct{}{}
	}
	return seen, nil
}
