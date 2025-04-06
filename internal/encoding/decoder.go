package encoding

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/junpeng.ong/blog/internal/encoding/file/codecV1.go"
	"github.com/junpeng.ong/blog/internal/encoding/file/utils"
	"github.com/junpeng.ong/blog/internal/filepb"
)

const (
	fileMarker     = "BLOC"
	fileMarkerSize = 4

	footerLengthByteSize      = 4
	minimumMetadataBufferSize = 64
)

var (
	// file marker errors

	ErrorProblemDecodingFileMarker = errors.New("problem decoding file marker")
	ErrorMissingInitialFileMarker  = errors.New("missing initial file marker")
	ErrorMissingFinalFileMarker    = errors.New("missing final file marker")

	// file metadata errors

	ErrorProblemDecodingFileVersion      = errors.New("problem decoding file version")
	ErrorProblemDecodingFileMetadata     = errors.New("problem decoding file metadata")
	ErrorProblemDecodingFileMetadataSize = errors.New("problem decoding file metadata size")
	ErrorMalformedFileMetadata           = errors.New("malformed file metadata")

	ErrorInvalidFileV1 = errors.New("invalid file v1")
)

type BlockFileLazyDecoder interface {
	Length() int
	DecodeContent(off int, size int) (string, error)
	DecodeBase64Image(off int, size int) (string, error)
	DecodeSection(index int) (*filepb.SectionNode, error)
}

func NewBlockFileLazyDecoder(reader io.ReadSeeker) (BlockFileLazyDecoder, error) {
	err := verifyFileMarkers(reader)
	if err != nil {
		return nil, err
	}

	footerStart, metadata, err := decodeMetadata(reader)
	if err != nil {
		return nil, err
	}

	switch metadata.Version {
	case 1:
		// validate content ranges
		cursor := fileMarkerSize
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
		file_metadata := metadata.GetFileMetadata()
		if cursor != footerStart {
			return nil, fmt.Errorf("%w: invalid metadata start offset: expected %d but cursor at %d", ErrorInvalidFileV1, footerStart, cursor)
		}
		// verify that the file cursor is
		cursor += int(file_metadata.SizeVT() + fileMarkerSize + footerLengthByteSize)
		if cursor != int(metadata.GetSize()) {
			return nil, fmt.Errorf("%w: expected file size of %d but got %d", ErrorInvalidFileV1, metadata.GetSize(), cursor)
		}

		return codecV1.NewBlockFileDecoderV1(reader, metadata), nil

	default:
		return nil, nil
	}
}

func verifyFileMarkers(reader io.ReadSeeker) error {
	var err error
	b := make([]byte, fileMarkerSize)
	// reading the starting byte marker
	_, err = utils.ReadFromStartOffset(reader, b, 0, fileMarkerSize)
	if err != nil {
		return fmt.Errorf("%w:%v", ErrorProblemDecodingFileMarker, err)
	}
	if bytes.Equal(b, []byte(fileMarker)) {
		// suggests that the file is likely not decodable by this codec
		return ErrorMissingInitialFileMarker
	}
	// reading the ending byte marker
	_, err = utils.ReadFromEndOffset(reader, b, fileMarkerSize, fileMarkerSize)
	if err != nil {
		return fmt.Errorf("%w:%v", ErrorProblemDecodingFileMarker, err)
	}
	if bytes.Equal(b, []byte(fileMarker)) {
		// suggests that the file is malformed
		// if the file was read from a stream, then the file may be incomplete
		return ErrorMissingFinalFileMarker
	}
	return nil
}

func decodeMetadata(reader io.ReadSeeker) (int, *filepb.Metadata, error) {
	var err error
	b := make([]byte, minimumMetadataBufferSize)
	// reading the footer size
	_, err = utils.ReadFromEndOffset(reader, b, fileMarkerSize+footerLengthByteSize, footerLengthByteSize)
	if err != nil {
		return 0, nil, fmt.Errorf("%w:%v", ErrorProblemDecodingFileMetadataSize, err)
	}
	footerSize := int(binary.LittleEndian.Uint32(b[:footerLengthByteSize]))

	// reading the footer
	if footerSize > minimumMetadataBufferSize {
		b = make([]byte, footerSize)
	}
	footerStart := footerSize + fileMarkerSize + footerLengthByteSize
	_, err = utils.ReadFromEndOffset(reader, b, footerStart, footerSize)
	if err != nil {
		return 0, nil, fmt.Errorf("%w:%v", ErrorProblemDecodingFileMetadata, err)
	}

	var metadata filepb.Metadata
	if err := metadata.UnmarshalVT(b); err != nil {
		return 0, nil, fmt.Errorf("%w: %v", ErrorMalformedFileMetadata, err)
	}
	return footerStart, &metadata, nil
}
