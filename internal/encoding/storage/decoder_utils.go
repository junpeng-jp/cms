package storage

import (
	"encoding/binary"
	"fmt"

	"github.com/junpeng.ong/blog/internal/file/contentpb"
	"github.com/junpeng.ong/blog/internal/file/metadatapb"
)

func StripFileMarkers(b []byte) ([]byte, error) {
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

func ComputeFooterByteRange(src []byte) (int, int) {
	// read 4 bytes that records the length of the footer
	leftOffset := len(src) - footerLengthByteSize
	footerSize := int(binary.LittleEndian.Uint32(src[leftOffset:]))

	return leftOffset - footerSize, leftOffset
}

func DecodeFooter(src []byte) (*metadatapb.Footer, error) {
	var footer metadatapb.Footer
	if err := footer.UnmarshalVT(src); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrorMalformedFileFooter, err)
	}
	return &footer, nil
}

func DecodeContentList(src []byte, footer *metadatapb.Footer) (int, []*contentpb.Content, error) {
	contentList := make([]*contentpb.Content, len(footer.ContentRange))

	var leftOffset, rightOffset int
	for i, contentMeta := range footer.ContentRange {
		if int(contentMeta.Offset) != leftOffset {
			return 0, nil, fmt.Errorf("%w: expected to start at %d but current offset is at %d", ErrorMalformedContent, contentMeta.Offset, leftOffset)
		}
		rightOffset = leftOffset + int(contentMeta.Size)

		var content contentpb.Content
		if err := content.UnmarshalVT(src[leftOffset:rightOffset]); err != nil {
			return 0, nil, fmt.Errorf("%w: %v", ErrorMalformedContent, err)
		}
		contentList[i] = &content

		leftOffset = rightOffset
	}

	return rightOffset, contentList, nil
}
