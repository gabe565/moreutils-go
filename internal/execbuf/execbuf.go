package execbuf

import (
	"errors"
	"io"
	"os/exec"
	"slices"
	"sync"
	"time"

	"github.com/gabe565/moreutils/internal/util"
)

// Buffer is an io.Writer that buffers exec.Cmd output.
// Writes contain timestamp and source metadata so the output can be replayed.
type Buffer struct {
	writes []write
	mu     sync.Mutex
	wg     sync.WaitGroup
}

var (
	ErrStdoutSet = errors.New("stdout already set")
	ErrStderrSet = errors.New("stderr already set")
	ErrStarted   = errors.New("exec buffer after process started")
)

// RunBuffered runs the provided exec.Cmd with buffered stdout/stderr streams.
func RunBuffered(cmd *exec.Cmd, stdout, stderr io.Writer) (*Buffer, error) {
	if cmd.Stdout != nil {
		return nil, ErrStdoutSet
	}
	if cmd.Stderr != nil {
		return nil, ErrStderrSet
	}
	if cmd.Process != nil {
		return nil, ErrStarted
	}

	buf := &Buffer{}
	cmd.Stdout = buf.Writer(stdout)
	cmd.Stderr = buf.Writer(stderr)
	err := cmd.Run()
	return buf, err
}

// write represents a single write
type write struct {
	ts     time.Time
	source io.Writer
	data   []byte
}

// Writer creates an BufferWriter which writes to the provided stream
func (e *Buffer) Writer(f io.Writer) *BufferWriter {
	return &BufferWriter{
		Buffer: e,
		source: f,
	}
}

// Close waits for all in-flight writes to finish, then sorts them by timestamp
func (e *Buffer) Close() {
	e.wg.Wait()

	e.mu.Lock()
	defer e.mu.Unlock()

	slices.SortStableFunc(e.writes, func(a, b write) int {
		return a.ts.Compare(b.ts)
	})
}

// Print prints output for a source, or all sources if nil.
func (e *Buffer) Print(source io.Writer) error {
	e.Close()
	e.mu.Lock()
	defer e.mu.Unlock()

	var errs []error
	for _, write := range e.writes {
		if source == nil || write.source == source {
			if _, err := write.source.Write(write.data); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return util.JoinErrors(errs...)
}

// Bytes returns the bytes for a source, or all sources if nil.
func (e *Buffer) Bytes(source io.Writer) []byte {
	e.Close()
	buf := make([]byte, 0, e.Len(source))
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, write := range e.writes {
		if source == nil || write.source == source {
			buf = append(buf, write.data...)
		}
	}
	return buf
}

// Len returns the number of bytes written for a source, or all sources if nil.
func (e *Buffer) Len(source io.Writer) int64 {
	e.mu.Lock()
	defer e.mu.Unlock()

	var n int64
	for _, write := range e.writes {
		if source == nil || write.source == source {
			n += int64(len(write.data))
		}
	}
	return n
}

// BufferWriter handles writes to a given source
type BufferWriter struct {
	*Buffer
	source io.Writer
}

// Write writes bytes for a source to the parent Buffer
func (e *BufferWriter) Write(p []byte) (int, error) {
	ts := time.Now()
	data := slices.Clone(p)
	e.wg.Add(1)

	go func() {
		defer e.wg.Done()

		e.mu.Lock()
		defer e.mu.Unlock()

		e.writes = append(e.writes, write{
			ts:     ts,
			source: e.source,
			data:   data,
		})
	}()

	return len(p), nil
}
