package testutils

import (
	"errors"
	"io"
)

type SeekableWriter struct {
	Buf []byte
	pos int
}

func (m *SeekableWriter) Write(p []byte) (n int, err error) {
	if len(p) > cap(m.Buf)-m.pos {
		m.Buf = append(m.Buf, make([]byte, len(p))...)[:m.pos]
	}
	m.Buf = m.Buf[:m.pos+len(p)]
	n = copy(m.Buf[m.pos:], p)
	m.pos += n
	return n, nil
}

func (m *SeekableWriter) Seek(offset int64, whence int) (int64, error) {
	newPos, offs := 0, int(offset)
	switch whence {
	case io.SeekStart:
		newPos = offs
	case io.SeekCurrent:
		newPos = m.pos + offs
	case io.SeekEnd:
		newPos = len(m.Buf) + offs
	}
	if newPos < 0 {
		return 0, errors.New("negative result pos")
	}
	m.pos = newPos
	return int64(newPos), nil
}
