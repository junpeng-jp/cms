package encoding

import (
	"github.com/junpeng.ong/blog/internal/encoding/storage"
	"github.com/junpeng.ong/blog/internal/file"
)

var (
	// Codec to convert to and from the binary storage format to a File
	_ Encoder = &storage.StorageCodec{}
	_ Decoder = &storage.StorageCodec{}
)

type Decoder interface {
	Decode([]byte) (*file.File, error)
}

type Encoder interface {
	Encode(*file.File) ([]byte, error)
}
