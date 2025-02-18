package combine

import (
	"errors"
	"fmt"
	"io"
)

//go:generate go tool enumer -type operator -trimprefix operator -transform lower -output operator_string.go

type operator uint8

const (
	operatorAnd operator = iota
	operatorNot
	operatorOr
	operatorXor
)

var ErrInvalidOperator = errors.New("invalid operator")

// compare runs the correct compare function for the current operator.
func (op operator) compare(out io.Writer, r1, r2 io.ReadSeeker) error {
	switch op {
	case operatorAnd:
		return compareAnd(out, r1, r2)
	case operatorNot:
		return compareNot(out, r1, r2)
	case operatorOr:
		return compareOr(out, r1, r2)
	case operatorXor:
		return compareXor(out, r1, r2)
	default:
		return ErrInvalidOperator
	}
}

// compareOr outputs lines from both r1 and r2
func compareOr(out io.Writer, r1, r2 io.Reader) error {
	for line, err := range iterLines(r1) {
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(out, line); err != nil {
			return err
		}
	}
	for line, err := range iterLines(r2) {
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(out, line); err != nil {
			return err
		}
	}
	return nil
}

// compareXor outputs lines that are in r1 or r2, but not in both
func compareXor(out io.Writer, r1, r2 io.ReadSeeker) error {
	if err := compareNot(out, r1, r2); err != nil {
		return err
	}

	if _, err := r1.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if _, err := r2.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if err := compareNot(out, r2, r1); err != nil {
		return err
	}

	return nil
}

// compareNot outputs lines from r1 that are not in r2
func compareNot(out io.Writer, r1, r2 io.Reader) error {
	seen, err := collectLines(r2)
	if err != nil {
		return err
	}

	for line, err := range iterLines(r1) {
		if err != nil {
			return err
		}
		if _, exists := seen[line]; !exists {
			if _, err := fmt.Fprintln(out, line); err != nil {
				return err
			}
		}
	}
	return nil
}

// compareAnd outputs lines that are in both r1 and r2
func compareAnd(out io.Writer, r1, r2 io.Reader) error {
	seen, err := collectLines(r2)
	if err != nil {
		return err
	}

	for line, err := range iterLines(r1) {
		if err != nil {
			return err
		}
		if _, exists := seen[line]; exists {
			if _, err := fmt.Fprintln(out, line); err != nil {
				return err
			}
		}
	}
	return nil
}
