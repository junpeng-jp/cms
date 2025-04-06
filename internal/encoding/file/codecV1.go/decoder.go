package codecV1

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/junpeng.ong/blog/internal/encoding/utils"
	"github.com/junpeng.ong/blog/internal/filepb"
)

const v1InitialDecoderBufferSize = 1024 // 1KB

var (
	// section decoding errors

	ErrorSectionIndexOutOfRange    = errors.New("section index out of range")
	ErrorFailedToDecodeFileSection = errors.New("failed to decode file section")
	ErrorMalformedFileSection      = errors.New("malformed file section")

	// content decoding errors

	ErrorFailedToDecodeFileContent = errors.New("failed to decode file content")
	ErrorContentIsNotValidUTF8     = errors.New("content is not valid utf8")

	// image decoding errors

	ErrorFailedToDecodeFileImage     = errors.New("failed to decode file image")
	ErrorImageIsNotValidBase64String = errors.New("image is not valid base64 string")

	ErrorBufferTooLarge = errors.New("buffer too large")
)

type blockFileDecoderV1 struct {
	reader   io.ReadSeeker
	metadata *filepb.Metadata
	buffer   []byte
}

func NewBlockFileDecoderV1(reader io.ReadSeeker, metadata *filepb.Metadata) *blockFileDecoderV1 {
	return &blockFileDecoderV1{
		reader:   reader,
		metadata: metadata,
		buffer:   make([]byte, v1InitialDecoderBufferSize),
	}
}

func (d *blockFileDecoderV1) Length() int {
	return len(d.metadata.SectionMetadata.Ranges)
}

func (d *blockFileDecoderV1) DecodeContent(off int, size int) (string, error) {
	var err error
	if err := d.resizeBuffer(size); err != nil {
		return "", err
	}

	b := d.buffer[:size]
	_, err = utils.ReadFromStartOffset(d.reader, b, off, size)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrorFailedToDecodeFileContent, err)
	}

	if !utf8.Valid(b) {
		return "", ErrorContentIsNotValidUTF8
	}
	return string(b), nil
}

func (d *blockFileDecoderV1) DecodeBase64Image(off int, size int) (string, error) {
	var err error
	if err := d.resizeBuffer(size); err != nil {
		return "", err
	}

	b := d.buffer[:size]
	_, err = utils.ReadFromStartOffset(d.reader, b, off, size)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrorFailedToDecodeFileImage, err)
	}

	if !base64UrlSafeRegex.Match(b) {
		return "", fmt.Errorf("%w: %v", ErrorImageIsNotValidBase64String, err)
	}

	return string(b), nil
}

func (d *blockFileDecoderV1) DecodeSection(index int) (*filepb.SectionNode, error) {
	if index >= d.Length() {
		return nil, ErrorSectionIndexOutOfRange
	}

	section := d.metadata.SectionMetadata.Ranges[index]
	off := int(section.Start)
	size := int(section.End - section.Start)

	var err error
	if err = d.resizeBuffer(size); err != nil {
		return nil, err
	}

	b := d.buffer[:size]
	_, err = utils.ReadFromStartOffset(d.reader, b, off, size)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrorFailedToDecodeFileSection, err)
	}

	var sectionNode filepb.SectionNode
	if err = sectionNode.UnmarshalVT(b); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrorMalformedFileSection, err)
	}
	return &sectionNode, nil
}

func (d *blockFileDecoderV1) resizeBuffer(n int) error {
	if n > v1MaxContentTotalSize {
		return ErrorBufferTooLarge
	}
	if len(d.buffer) < n {
		d.buffer = make([]byte, n)
	}
	return nil
}
