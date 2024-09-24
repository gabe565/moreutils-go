//go:build !unix

package util

func Umask(_ int) int {
	return 0
}
