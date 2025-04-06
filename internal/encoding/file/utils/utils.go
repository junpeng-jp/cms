package utils

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrorInsufficientBytesRead    = errors.New("insufficient bytes read")
	ErrorInsufficientBytesWritten = errors.New("insufficient bytes written")
)

func ReadFromStartOffset(reader io.ReadSeeker, b []byte, offset int, size int) (n int, err error) {
	_, err = reader.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return 0, err
	}
	n, err = reader.Read(b[:size])
	if n != size {
		return n, fmt.Errorf("%w: expected to read %db but only read %db", ErrorInsufficientBytesRead, size, n)
	}
	return n, err
}

func ReadFromEndOffset(reader io.ReadSeeker, b []byte, offset int, size int) (n int, err error) {
	_, err = reader.Seek(int64(offset), io.SeekEnd)
	if err != nil {
		return 0, err
	}
	n, err = reader.Read(b[:size])
	if n != size {
		return n, fmt.Errorf("%w: expected to read %db but only read %db", ErrorInsufficientBytesRead, size, n)
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
