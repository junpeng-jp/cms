package encoding

import (
	"io"

	"github.com/junpeng.ong/blog/internal/encoding/file/codecV1.go"
	"github.com/junpeng.ong/blog/internal/filepb"
)

type BlockFileEncoder interface {
	GetFinalContentMetadata() *filepb.ByteRange
	GetFinalSectionMetadata() *filepb.SectionMetadata
	EncodeSectionImage(*filepb.SectionNode, []byte) error
	EncodeSectionContent(*filepb.SectionNode, []byte) error
	Finalize() (int, error)
}

func NewBlockFileEncoder(writer io.WriteSeeker) (BlockFileEncoder, error) {
	return codecV1.NewBlockFileEncoderV1(writer, fileMarkerSize), nil
}
