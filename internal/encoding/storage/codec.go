package storage

import (
	"github.com/junpeng.ong/blog/internal/file"
	"github.com/junpeng.ong/blog/internal/file/filepb"
)
type ContentDecoder func([]byte) error

type ContentEncoder func() ([]byte, error)

type StorageCodec struct{}

func (p *StorageCodec) Decode(b []byte) (*file.File, error) {
	bNoMarker, err := StripFileMarkers(b)
	if err != nil {
		return nil, err
	}
	var f filepb.File
	if err := f.UnmarshalVT(bNoMarker); err != nil {
		return nil, err
	}

	f.

	return
}

func (p *StorageCodec) Encode(file *file.File) ([]byte, error) {
	
}
