package testutils

import (
	"strings"

	"github.com/junpeng.ong/blog/internal/encoding"
	"github.com/junpeng.ong/blog/internal/filepb"
)

const letters = "abcdefghijklmnopqrstuvwxyz=-+/ "

var (
	_ encoding.BlockFileLazyDecoder = &mockDecoder{}
)

type mockDecoder struct {
	maxSize int
}

func NewMockDecoder(size int) *mockDecoder {
	return &mockDecoder{maxSize: size}
}

func (d *mockDecoder) Length() int {
	return d.maxSize
}

func (d *mockDecoder) DecodeContent(off int, size int) (string, error) {
	var writer strings.Builder
	for {
		if size <= 0 {
			break
		}
		writer.WriteString(letters[:min(size, len(letters))])
		size = max(size-len(letters), 0)
	}

	return writer.String(), nil
}

func (d *mockDecoder) DecodeBase64Image(off int, size int) (string, error) {
	var writer strings.Builder
	for {
		if size == 0 {
			break
		}
		writer.WriteString(letters[:min(size, len(letters))])
		size = max(size-len(letters), 0)
	}

	return writer.String(), nil
}

func (d *mockDecoder) DecodeSection(index int) (*filepb.SectionNode, error) {
	return nil, nil
}
