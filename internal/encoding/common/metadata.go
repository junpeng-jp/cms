package common

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/junpeng.ong/blog/internal/encoding/utils"
	"github.com/junpeng.ong/blog/internal/filepb"
)

const (
	FooterLengthByteSize      = 4
	minimumMetadataBufferSize = 64
)

var (
	ErrorProblemDecodingFileVersion      = errors.New("problem decoding file version")
	ErrorProblemDecodingFileMetadata     = errors.New("problem decoding file metadata")
	ErrorProblemDecodingFileMetadataSize = errors.New("problem decoding file metadata size")
	ErrorMalformedFileMetadata           = errors.New("malformed file metadata")
)

func DecodeMetadata(reader io.ReadSeeker) (int, *filepb.Metadata, error) {
	var err error
	b := make([]byte, minimumMetadataBufferSize)
	// reading the footer size
	_, err = utils.ReadFromEndOffset(reader, b, -(FileMarkerSize + FooterLengthByteSize), FooterLengthByteSize)
	if err != nil {
		return 0, nil, fmt.Errorf("%w:%v", ErrorProblemDecodingFileMetadataSize, err)
	}
	footerSize := int(binary.LittleEndian.Uint32(b[:FooterLengthByteSize]))

	// reading the footer
	if footerSize > minimumMetadataBufferSize {
		b = make([]byte, footerSize)
	}
	footerStart := footerSize + FileMarkerSize + FooterLengthByteSize
	_, err = utils.ReadFromEndOffset(reader, b, -(footerStart), footerSize)
	if err != nil {
		return 0, nil, fmt.Errorf("%w:%v", ErrorProblemDecodingFileMetadata, err)
	}
	fmt.Printf("%d\n", footerStart)
	fmt.Printf("%s\n", b)

	var metadata filepb.Metadata
	if err := metadata.UnmarshalVT(b[:footerSize]); err != nil {
		return 0, nil, fmt.Errorf("%w: %v", ErrorMalformedFileMetadata, err)
	}
	return footerStart, &metadata, nil
}
