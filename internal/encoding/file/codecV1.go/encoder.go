package codecV1

import (
	"errors"
	"fmt"
	"io"

	"github.com/junpeng.ong/blog/internal/encoding/file/utils"
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
	writer        io.WriteSeeker
	contentStart  int
	offset        int
	finalized     bool
	sections      []*filepb.SectionNode
	sectionRanges []*filepb.ByteRange
}

func NewBlockFileEncoderV1(writer io.WriteSeeker, initialOffset int) *blockFileEncoderV1 {
	return &blockFileEncoderV1{
		writer:        writer,
		contentStart:  initialOffset,
		offset:        initialOffset,
		finalized:     false,
		sections:      nil,
		sectionRanges: nil,
	}
}

func (e *blockFileEncoderV1) GetFinalContentMetadata() *filepb.ByteRange {
	if !e.finalized {
		return nil
	}

	return &filepb.ByteRange{
		Start: int32(e.contentStart),
		End:   int32(e.sectionRanges[0].Start),
	}
}

func (e *blockFileEncoderV1) GetFinalSectionMetadata() *filepb.SectionMetadata {
	if !e.finalized {
		return nil
	}

	return &filepb.SectionMetadata{Ranges: e.sectionRanges}
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
	_, err = utils.WriteFromStartOffset(e.writer, base64Bytes, e.offset, n)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorFailedToWriteContent, err)
	}
	e.sections = append(e.sections, section)
	e.offset += n

	return nil
}

func (e *blockFileEncoderV1) EncodeSectionContent(section *filepb.SectionNode, content []byte) error {
	var err error
	if err = e.checkContentBeforeWrite(len(content)); err != nil {
		return err
	}
	n := len(content)
	_, err = utils.WriteFromStartOffset(e.writer, content, e.offset, n)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorFailedToWriteContent, err)
	}
	e.sections = append(e.sections, section)
	e.offset += n

	return nil
}

func (e *blockFileEncoderV1) checkContentBeforeWrite(contentSize int) error {
	if e.finalized {
		return ErrorAttemptToEncodeAfterEncoderHasFinalized
	}
	if contentSize > v1MaxContentTotalSize-e.offset {
		return fmt.Errorf("%w: max content size is %db: attempted to write %db + %db", ErrorContentTooLarge, v1MaxContentTotalSize, e.offset, contentSize)

	}
	if contentSize > v1MaxContentChunkSize {
		return fmt.Errorf("%w: max content chunk size is %db: got a chunk of size %db", ErrorContentChunkTooLarge, v1MaxContentTotalSize, contentSize)
	}
	return nil
}

func (e *blockFileEncoderV1) Finalize() (int, error) {

	buffer := make([]byte, v1InitialEncoderBufferSize)

	offset := e.offset

	if len(e.sections) > 0 {
		var b []byte
		var err error
		var n, size int

		sectionRanges := make([]*filepb.ByteRange, len(e.sections))

		for i, section := range e.sections {
			size = section.SizeVT()
			if len(buffer) < size {
				buffer = make([]byte, size)
			}
			b = buffer[:size]
			n, err = section.MarshalToSizedBufferVT(b)
			if err != nil {
				return 0, fmt.Errorf("%w: %v", ErrorFailedToEncodeSection, err)
			}
			if n != size {
				return 0, ErrorPartialEncodingOfSection
			}
			_, err = utils.WriteFromStartOffset(e.writer, b, offset, size)
			if err != nil {
				return 0, fmt.Errorf("%w: %v", ErrorFailedToWriteSection, err)
			}
			sectionRanges[i] = &filepb.ByteRange{
				Start: int32(offset),
				End:   int32(offset + size),
			}
			offset += size
		}

		e.sectionRanges = sectionRanges
	}

	e.offset = offset
	e.finalized = true

	return offset - e.contentStart, nil
}
