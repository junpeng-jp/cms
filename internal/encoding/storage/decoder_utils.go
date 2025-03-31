package storage

import (
	"encoding/binary"
	"fmt"

	"github.com/junpeng.ong/blog/internal/file"
	"github.com/junpeng.ong/blog/internal/file/filepb"
)

func DecodeFile(b []byte) (*file.File, error) {
	bNoMarker, err := StripFileMarkers(b)
	if err != nil {
		return nil, err
	}

	footerStart, footerEnd := ComputeFooterByteRange(bNoMarker)

	// read the footer to get file metadata
	footer, err := DecodeFooter(bNoMarker[footerStart:footerEnd])
	if err != nil {
		return nil, err
	}

	// read the content portion of the file
	n, contentList, err := DecodeContentList(bNoMarker[footer.ContentStartOffset:], footer)
	if err != nil {
		return nil, err
	}

	// read the document portion of the file
	if int(footer.DocumentStartOffset) != int(footer.ContentStartOffset)+n {
		// the left offset should be at the start of the document section
		return nil, nil
	}

	return &file.File{
		ContentList: contentList,
	}, nil
}

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

func DecodeFooter(src []byte) (*filepb.Footer, error) {
	var footer filepb.Footer
	if err := footer.UnmarshalVT(src); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrorMalformedFileFooter, err)
	}
	return &footer, nil
}

func DecodeContentList(src []byte, footer *filepb.Footer) (int, []*filepb.Content, error) {
	contentList := make([]*filepb.Content, len(footer.ContentRange))

	var leftOffset, rightOffset int
	for i, contentMeta := range footer.ContentRange {
		if int(contentMeta.Offset) != leftOffset {
			return 0, nil, fmt.Errorf("%w: expected to start at %d but current offset is at %d", ErrorMalformedContent, contentMeta.Offset, leftOffset)
		}
		rightOffset = leftOffset + int(contentMeta.Size)

		var content filepb.Content
		if err := content.UnmarshalVT(src[leftOffset:rightOffset]); err != nil {
			return 0, nil, fmt.Errorf("%w: %v", ErrorMalformedContent, err)
		}
		contentList[i] = &content

		content.Kind

		leftOffset = rightOffset
	}

	return rightOffset, contentList, nil
}

type ContentTransform struct {
	InlineDecoder func()	
}

func Decode