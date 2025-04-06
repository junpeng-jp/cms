package encoding

import (
	"io"

	"github.com/junpeng.ong/blog/internal/encoding/file/codecV1.go"
	"github.com/junpeng.ong/blog/internal/filepb"
)

type BlockFileEncoder interface {
	Init() error
	EncodeSectionImage(*filepb.SectionNode, []byte) error
	EncodeSectionContent(*filepb.SectionNode, []byte) error
	Finalize(string) (int, error)
}

func NewBlockFileEncoder(writer io.Writer) (BlockFileEncoder, error) {
	return codecV1.NewBlockFileEncoderV1(writer), nil
}
