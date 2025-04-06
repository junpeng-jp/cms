package utils

import (
	"errors"
	"fmt"
	"io"
)

type SeekableWriter struct {
	buf []byte
	pos int
}

func NewSeekableWriter(buf []byte) *SeekableWriter {
	return &SeekableWriter{buf: buf, pos: 0}
}

func (m *SeekableWriter) GetBuffer() []byte {
	return m.buf[:m.pos]
}

func (m *SeekableWriter) Write(p []byte) (n int, err error) {
	if len(p) > cap(m.buf)-m.pos {
		m.buf = append(m.buf, make([]byte, len(p))...)[:m.pos]
	}
	m.buf = m.buf[:m.pos+len(p)]
	n = copy(m.buf[m.pos:], p)
	m.pos += n
	return n, nil
}

func (m *SeekableWriter) Seek(offset int64, whence int) (int64, error) {
	newPos := 0
	switch whence {
	case io.SeekStart:
		newPos = int(offset)
	case io.SeekCurrent:
		newPos = m.pos + int(offset)
	case io.SeekEnd:
		newPos = len(m.buf) + int(offset)
	}
	if newPos < 0 {
		return 0, errors.New("negative result pos")
	}
	m.pos = newPos
	return int64(newPos), nil
}

func WriteFromCurrentPosition(writer io.Writer, b []byte, size int) (n int, err error) {
	n, err = writer.Write(b[:size])
	if n != size {
		return n, fmt.Errorf("%w: expected to write %db but only wrote %db", ErrorInsufficientBytesWritten, size, n)
	}
	return n, err
}

func WriteFromStartOffset(writer io.WriteSeeker, b []byte, offset int, size int) (n int, err error) {
	_, err = writer.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return 0, err
	}
	n, err = writer.Write(b[:size])
	if n != size {
		return n, fmt.Errorf("%w: expected to write %db but only wrote %db", ErrorInsufficientBytesWritten, size, n)
	}
	return n, err
}

func WriteFromEndOffset(writer io.WriteSeeker, b []byte, offset int, size int) (n int, err error) {
	_, err = writer.Seek(int64(offset), io.SeekEnd)
	if err != nil {
		return 0, err
	}
	n, err = writer.Write(b[:size])
	if n != size {
		return n, fmt.Errorf("%w: expected to write %db but only wrote %db", ErrorInsufficientBytesWritten, size, n)
	}
	return n, err
}
