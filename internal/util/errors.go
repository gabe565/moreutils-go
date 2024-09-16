package util

import (
	"errors"
	"slices"
)

// JoinErrors behaves similarly to errors.Join, but returns the error verbatim if there is only 1.
func JoinErrors(errs ...error) error {
	errs = slices.DeleteFunc(errs, func(err error) bool {
		return err == nil
	})
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		return errors.Join(errs...)
	}
}
