//go:build unix

package errno

import "iter"

func Iter() iter.Seq[*Errno] {
	return func(yield func(*Errno) bool) {
		for num := range 256 {
			if errno := New(num); errno.Valid() {
				if !yield(errno) {
					return
				}
			}
		}

		// MIPS has an error code > 256
		if errno := New(1133); errno.Valid() {
			if !yield(errno) {
				return
			}
		}
	}
}
