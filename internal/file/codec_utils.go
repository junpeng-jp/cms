package file

import (
	"encoding/binary"
	"fmt"

	"github.com/junpeng.ong/blog/internal/file/filepb"
)

const (
	footerLengthByteSize = 4

	fileMarker     = "BLOC"
	fileMarkerSize = 4

	maxFileSize = 4194304 // 4MB
)

func stripFileMarkers(b []byte) ([]byte, error) {
	if string(b[:fileMarkerSize]) != fileMarker {
		// suggests that the file is likely not decodable by this codec
		return nil, ErrorMissingInitialFileMarker
	}
	leftOffset := len(b) - fileMarkerSize

	if string(b[leftOffset:]) != fileMarker {
		// suggests that the file is malformed
		// if the file was read from a stream, then the file may be incomplete
		return nil, ErrorMissingFinalFileMarker
	}

	return b[fileMarkerSize:leftOffset], nil
}

func computeMetadataByteRange(src []byte) (int, int) {
	// read 4 bytes that records the length of the footer
	leftOffset := len(src) - footerLengthByteSize
	footerSize := int(binary.LittleEndian.Uint32(src[leftOffset:]))

	return leftOffset - footerSize, leftOffset
}

func decodeMetadata(src []byte) (*filepb.Metadata, error) {
	var metadata filepb.Metadata
	if err := metadata.UnmarshalVT(src); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrorMalformedFileMetadata, err)
	}
	return &metadata, nil
}
