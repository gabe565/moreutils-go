//go:build unix

package errno

import "iter"

func Iter() iter.Seq[*Errno] {
	return func(yield func(*Errno) bool) {
		for num := range 2048 {
			if errno := New(num); errno.Valid() {
				if !yield(errno) {
					return
				}
			}
		}
	}
}
