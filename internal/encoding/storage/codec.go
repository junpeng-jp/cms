package storage

import (
	"encoding/binary"
	"fmt"

	"github.com/junpeng.ong/blog/internal/file"
	"github.com/junpeng.ong/blog/internal/file/filepb"
)

const (
	footerVersion        = 1
	footerLengthByteSize = 4

	fileMarker     = "BLOC"
	fileMarkerSize = 4

	maxFileSize = 4194304 // 4MB
)

type ContentDecoder func([]byte) error

type ContentEncoder func() ([]byte, error)

type StorageCodec struct{}

func (p *StorageCodec) Decode(b []byte) (*file.File, error) {
	bNoMarker, err := StripFileMarkers(b)
	if err != nil {
		return nil, err
	}
	var f filepb.File
	if err := f.UnmarshalVT(bNoMarker); err != nil {
		return nil, err
	}

	return
}

func (p *StorageCodec) Encode(file *file.File) ([]byte, error) {
	contentList := file.ContentList
	contentRange := make([]*filepb.ContentRange, len(contentList))

	var leftOffset int
	var rightOffset int
	var err error

	for i, content := range contentList {
		n := content.SizeVT()
		rightOffset += n
		if rightOffset > maxFileSize {
			return nil, fmt.Errorf("%w: exceeded maximum length of %d by %d when encoding", ErrorFileExceedMaxSize, maxFileSize, rightOffset-maxFileSize)
		}
		contentRange[i] = &filepb.ContentRange{
			Length: 0,
			Size:   int32(n),
			Offset: int32(leftOffset),
		}
		leftOffset = rightOffset
	}

	footer := &filepb.Footer{
		Version:             footerVersion,
		ContentSize:         int32(rightOffset),
		ContentStartOffset:  0,
		ContentRange:        contentRange,
		DocumentSize:        0,
		DocumentStartOffset: int32(rightOffset),
	}

	footerSize := footer.SizeVT()

	totalSize := fileMarkerSize + rightOffset + footerSize + footerLengthByteSize + fileMarkerSize
	fileBytes := make([]byte, totalSize)

	var n int
	n = copy(fileBytes[:fileMarkerSize], fileMarker)
	if n != fileMarkerSize {
		return nil, fmt.Errorf("%w: expected %d: got %d", ErrorUnexpectedHeadFileMarkerSize, fileMarkerSize, n)
	}

	leftOffset = n
	rightOffset = n
	for i, contentMeta := range contentRange {
		if int(contentMeta.Offset)+fileMarkerSize != leftOffset {
			return nil, ErrorUnexpectedOffsetWhenEncodingFile
		}
		rightOffset = leftOffset + int(contentMeta.Size)
		_, err = contentList[i].MarshalToSizedBufferVT(fileBytes[leftOffset:rightOffset])
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrorUnableToEncodeFile, err)
		}
		leftOffset = rightOffset
	}

	rightOffset = rightOffset + footerSize
	footer.MarshalToSizedBufferVT(fileBytes[leftOffset:rightOffset])

	leftOffset = rightOffset
	rightOffset = rightOffset + footerLengthByteSize
	binary.LittleEndian.PutUint32(fileBytes[leftOffset:rightOffset], uint32(footerSize))

	leftOffset = rightOffset
	rightOffset = rightOffset + fileMarkerSize
	n = copy(fileBytes[leftOffset:], fileMarker)
	if n != fileMarkerSize {
		return nil, fmt.Errorf("%w: expected %d: got %d", ErrorUnexpectedTailFileMarkerSize, fileMarkerSize, n)
	}

	return fileBytes[:rightOffset], nil
}
