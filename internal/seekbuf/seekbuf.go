package seekbuf

import "io"

// Buffer is an io.ReadCloser that supports seeking.
type Buffer struct {
	data []byte
	pos  int64
}

func New(data []byte) *Buffer {
	return &Buffer{
		data: data,
	}
}

func (b *Buffer) Read(p []byte) (int, error) {
	if b.pos >= int64(len(b.data)) {
		return 0, io.EOF
	}

	n := copy(p, b.data[b.pos:])
	b.pos += int64(n)
	return n, nil
}

// Seek sets the read pointer to offset.
func (b *Buffer) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		b.pos = offset
	case io.SeekCurrent:
		b.pos += offset
	case io.SeekEnd:
		b.pos = int64(len(b.data)) + offset
	}
	return b.pos, nil
}

// Close clears all the data out of the buffer and sets the read position to 0.
func (b *Buffer) Close() error {
	b.data = nil
	b.pos = 0
	return nil
}
