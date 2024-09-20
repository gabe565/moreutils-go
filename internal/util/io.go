package util

import (
	"io"
	"sync"
	"syscall"
)

func NewSuppressErrorWriter(w io.Writer) *SuppressErrorWriter {
	return &SuppressErrorWriter{w: w}
}

// SuppressErrorWriter proxies writes to another writer.
// If a write returns an error, the error will be suppressed and the writer will be disabled.
type SuppressErrorWriter struct {
	w      io.Writer
	closed bool
	mu     sync.Mutex
	Error  error
}

func (w *SuppressErrorWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	switch {
	case w.closed:
		if w.Error != nil {
			return 0, w.Error
		}
		return 0, syscall.EPIPE
	case w.Error == nil:
		_, w.Error = w.w.Write(p)
	}
	return len(p), nil
}

func (w *SuppressErrorWriter) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.closed = false
	w.Error = nil
}

func (w *SuppressErrorWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.closed = true
	if closer, ok := w.w.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
