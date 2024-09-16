package util

import "io"

func NewSuppressErrorWriter(w io.Writer) *SuppressErrorWriter {
	return &SuppressErrorWriter{w: w}
}

// SuppressErrorWriter proxies writes to another writer.
// If a write returns an error, the error will be suppressed and the writer will be disabled.
type SuppressErrorWriter struct {
	w     io.Writer
	Error error
}

func (w *SuppressErrorWriter) Write(p []byte) (int, error) {
	if w.Error == nil {
		_, w.Error = w.w.Write(p)
	}
	return len(p), nil
}

func (w *SuppressErrorWriter) Reset() {
	w.Error = nil
}
