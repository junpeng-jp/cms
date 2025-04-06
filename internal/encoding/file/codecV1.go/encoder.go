package codecV1

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/junpeng.ong/blog/internal/encoding/common"
	"github.com/junpeng.ong/blog/internal/encoding/utils"
	"github.com/junpeng.ong/blog/internal/filepb"
)

const v1InitialEncoderBufferSize = 1024 // 1KB

var (
	// content write errors

	ErrorAttemptToEncodeAfterEncoderHasFinalized = errors.New("attempt to encode after encoder has finalized")
	ErrorContentTooLarge                         = errors.New("content too large")
	ErrorContentChunkTooLarge                    = errors.New("content chunk large")
	ErrorFailedToWriteContent                    = errors.New("failed to write content")

	// section write errors

	ErrorFailedToEncodeSection    = errors.New("failed to encode section")
	ErrorPartialEncodingOfSection = errors.New("partial encoding of section")
	ErrorFailedToWriteSection     = errors.New("failed to write section")
)

type blockFileEncoderV1 struct {
	writer        io.Writer
	pos           int
	contentStart  int
	initialized   bool
	finalized     bool
	sections      []*filepb.SectionNode
	sectionRanges []*filepb.ByteRange
}

func NewBlockFileEncoderV1(writer io.Writer) *blockFileEncoderV1 {
	return &blockFileEncoderV1{
		writer:        writer,
		pos:           0,
		contentStart:  0,
		initialized:   false,
		finalized:     false,
		sections:      nil,
		sectionRanges: nil,
	}
}
func (e *blockFileEncoderV1) Init() error {
	_, err := utils.WriteFromCurrentPosition(e.writer, []byte(common.FileMarker), common.FileMarkerSize)
	if err != nil {
		return err
	}
	e.contentStart = common.FileMarkerSize
	e.pos = common.FileMarkerSize
	return nil
}

func (e *blockFileEncoderV1) EncodeSectionImage(section *filepb.SectionNode, base64Bytes []byte) error {
	var err error
	if err = e.checkContentBeforeWrite(len(base64Bytes)); err != nil {
		return err
	}

	if !base64UrlSafeRegex.Match(base64Bytes) {
		return fmt.Errorf("%w: %v", ErrorImageIsNotValidBase64String, err)
	}

	n := len(base64Bytes)
	_, err = utils.WriteFromCurrentPosition(e.writer, base64Bytes, n)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorFailedToWriteContent, err)
	}
	e.sections = append(e.sections, section)
	e.pos += n

	return nil
}

func (e *blockFileEncoderV1) EncodeSectionContent(section *filepb.SectionNode, content []byte) error {
	var err error
	if err = e.checkContentBeforeWrite(len(content)); err != nil {
		return err
	}
	n := len(content)
	_, err = utils.WriteFromCurrentPosition(e.writer, content, n)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorFailedToWriteContent, err)
	}
	e.sections = append(e.sections, section)
	e.pos += n

	return nil
}

func (e *blockFileEncoderV1) checkContentBeforeWrite(contentSize int) error {
	if e.finalized {
		return ErrorAttemptToEncodeAfterEncoderHasFinalized
	}
	if contentSize > v1MaxContentTotalSize-e.pos {
		return fmt.Errorf("%w: max content size is %db: attempted to write %db + %db", ErrorContentTooLarge, v1MaxContentTotalSize, e.pos, contentSize)

	}
	if contentSize > v1MaxContentChunkSize {
		return fmt.Errorf("%w: max content chunk size is %db: got a chunk of size %db", ErrorContentChunkTooLarge, v1MaxContentTotalSize, contentSize)
	}
	return nil
}

func (e *blockFileEncoderV1) Finalize(filename string) (int, error) {

	buffer := make([]byte, v1InitialEncoderBufferSize)

	pos := e.pos
	contentMetadata := &filepb.ByteRange{
		Start: int32(e.contentStart),
		End:   int32(pos),
	}

	var b []byte
	var err error
	var n, size int
	var sectionRanges []*filepb.ByteRange

	if len(e.sections) > 0 {
		sectionRanges = make([]*filepb.ByteRange, len(e.sections))

		for i, section := range e.sections {
			size = section.SizeVT()
			if len(buffer) < size {
				buffer = make([]byte, size)
			}
			b = buffer[:size]
			n, err = section.MarshalToSizedBufferVT(b)
			if err != nil {
				return pos + n, fmt.Errorf("%w: %v", ErrorFailedToEncodeSection, err)
			}
			if n != size {
				return pos + n, ErrorPartialEncodingOfSection
			}
			n, err = utils.WriteFromCurrentPosition(e.writer, b, size)
			if err != nil {
				return pos + n, fmt.Errorf("%w: %v", ErrorFailedToWriteSection, err)
			}
			sectionRanges[i] = &filepb.ByteRange{
				Start: int32(pos),
				End:   int32(pos + size),
			}
			pos += size
		}
	}
	// write metadata
	metadata := &filepb.Metadata{
		Version:         1,
		ContentMetadata: contentMetadata,
		SectionMetadata: &filepb.SectionMetadata{
			Ranges: sectionRanges,
		},
		FileMetadata: &filepb.FileMetadata{
			Name:      filename,
			CreatedAt: time.Now().UnixMilli(),
		},
	}
	size = metadata.SizeVT()
	if len(buffer) < size {
		buffer = make([]byte, size)
	}
	b = buffer[:size]
	n, err = metadata.MarshalToSizedBufferVT(b)
	if err != nil {
		return pos + n, fmt.Errorf("%w: %v", ErrorFailedToEncodeSection, err)
	}
	n, err = utils.WriteFromCurrentPosition(e.writer, b, size)
	if err != nil {
		return pos + n, fmt.Errorf("%w: %v", ErrorFailedToWriteSection, err)
	}
	pos += size

	// write metadata size
	b = buffer[:common.FooterLengthByteSize]
	binary.LittleEndian.PutUint32(b, uint32(size))

	n, err = utils.WriteFromCurrentPosition(e.writer, b, common.FooterLengthByteSize)
	if err != nil {
		return pos + n, fmt.Errorf("%w: %v", ErrorFailedToWriteSection, err)
	}

	pos += common.FooterLengthByteSize

	// write final file marker
	n, err = utils.WriteFromCurrentPosition(e.writer, []byte(common.FileMarker), common.FileMarkerSize)
	if err != nil {
		return pos + n, err
	}

	pos += common.FileMarkerSize

	e.pos = pos
	e.finalized = true

	return pos, nil
}
