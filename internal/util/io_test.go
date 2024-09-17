package util

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errorWriter struct {
	w io.Writer
}

var errFailed = errors.New("failed")

func (e errorWriter) Write(p []byte) (int, error) {
	n, _ := e.w.Write(p)
	return n, errFailed
}

func TestSuppressErrorWriter_Write(t *testing.T) {
	var buf strings.Builder
	w := NewSuppressErrorWriter(errorWriter{w: &buf})
	require.NotNil(t, w)

	// Initial bytes are written
	n, err := w.Write([]byte("test"))
	assert.Equal(t, 4, n)
	require.NoError(t, err)
	assert.Equal(t, "test", buf.String())

	// Further bytes are ignored
	n, err = w.Write([]byte("test"))
	assert.Equal(t, 4, n)
	require.NoError(t, err)
	assert.Equal(t, "test", buf.String())
}
