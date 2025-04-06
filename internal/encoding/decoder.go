package encoding

import (
	"errors"
	"fmt"
	"io"

	"github.com/junpeng.ong/blog/internal/encoding/common"
	"github.com/junpeng.ong/blog/internal/encoding/file/codecV1.go"
	"github.com/junpeng.ong/blog/internal/filepb"
)

var (
	ErrorInvalidFileV1 = errors.New("invalid file v1")
)

type BlockFileLazyDecoder interface {
	Length() int
	DecodeContent(off int, size int) (string, error)
	DecodeBase64Image(off int, size int) (string, error)
	DecodeSection(index int) (*filepb.SectionNode, error)
}

func NewBlockFileLazyDecoder(reader io.ReadSeeker) (BlockFileLazyDecoder, error) {
	err := common.VerifyFileMarkers(reader)
	if err != nil {
		return nil, err
	}

	footerStart, metadata, err := common.DecodeMetadata(reader)
	if err != nil {
		return nil, err
	}

	switch metadata.Version {
	case 1:
		// validate content ranges
		cursor := common.FileMarkerSize
		if cursor != int(metadata.GetContentMetadata().GetStart()) {
			return nil, fmt.Errorf("%w: invalid content start offset: expected %d but cursor at %d", ErrorInvalidFileV1, metadata.ContentMetadata.Start, cursor)
		}
		// validate section ranges
		cursor = int(metadata.GetContentMetadata().GetEnd())
		section_metadata := metadata.GetSectionMetadata()
		if len(section_metadata.GetRanges()) > 0 {
			// validate section ranges
			for i, sectionRange := range section_metadata.GetRanges() {
				if cursor != int(sectionRange.Start) {
					return nil, fmt.Errorf("%w: invalid section %d start offset: expected %d but cursor at %d", ErrorInvalidFileV1, i, sectionRange.Start, cursor)
				}
				cursor = int(sectionRange.End)
			}

		}
		// validate file metatdata ranges
		if cursor != footerStart {
			return nil, fmt.Errorf("%w: invalid metadata start offset: expected %d but cursor at %d", ErrorInvalidFileV1, footerStart, cursor)
		}

		return codecV1.NewBlockFileDecoderV1(reader, metadata), nil

	default:
		return nil, nil
	}
}
