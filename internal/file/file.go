package file

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"unicode/utf8"

	"github.com/junpeng.ong/blog/internal/file/filepb"
)

type BlockFile interface {
	Encode() ([]byte, error)
	DecodeContent(index int) (string, error)
	DecodeImage(index int) ([]byte, error)
	DecodeSection(index int) (*filepb.SectionNode, error)
}

func DecodeBlockFileFromBytes(src []byte) (BlockFile, error) {
	srcWithoutMarkers, err := stripFileMarkers(src)
	if err != nil {
		return nil, err
	}

	metadataStartOffset, metadataEndOffset := computeMetadataByteRange(srcWithoutMarkers)
	metadata, err := decodeMetadata(srcWithoutMarkers[metadataStartOffset:metadataEndOffset])
	if err != nil {
		return nil, err
	}

	switch metadata.Version {
	case 1:
		// extract content ranges
		var l int
		content := make([][]byte, len(metadata.ContentRanges))
		for i, contentRange := range metadata.ContentRanges {
			if int(contentRange.Start) != l {
				return nil, nil
			}
			content[i] = srcWithoutMarkers[contentRange.Start:contentRange.End]
			l = int(contentRange.End)
		}
		// extract image ranges
		images := make([][]byte, len(metadata.ImagesRanges))
		for i, imageRange := range metadata.ImagesRanges {
			if int(imageRange.Start) != l {
				return nil, nil
			}
			images[i] = srcWithoutMarkers[imageRange.Start:imageRange.End]
			l = int(imageRange.End)
		}
		// extract document ranges
		sections := make([][]byte, len(metadata.SectionRanges))
		for i, sectionRange := range metadata.SectionRanges {
			if int(sectionRange.Start) != l {
				return nil, nil
			}
			images[i] = srcWithoutMarkers[sectionRange.Start:sectionRange.End]
			l = int(sectionRange.End)
		}
		// verify that the offset has reached the metadata offset
		if l != metadataStartOffset {
			return nil, nil
		}

		return &blockFileV1{
			content:  content,
			images:   images,
			sections: sections,
		}, nil

	default:
		return nil, nil
	}
}

type blockFileV1 struct {
	content      [][]byte
	images       [][]byte
	sectionBytes [][]byte

	sections []*filepb.SectionNode
}

func (f *blockFileV1) Encode() ([]byte, error) {
	// heuristic to allocate an initial byte array size
	var estSize int
	for _, c := range f.content[:10] {
		estSize += len(c)
	}

	buffer := bytes.NewBuffer(make([]byte, estSize+2*fileMarkerSize))
	n, err := buffer.WriteString(fileMarker)
	if n != fileMarkerSize {
		return nil, fmt.Errorf("%w: expected %d: got %d", ErrorUnexpectedHeadFileMarkerSize, fileMarkerSize, n)
	}

	l := fileMarkerSize
	r := l
	var totalContentLength int
	contentRanges := make([]*filepb.ByteRange, len(f.content))
	for i, c := range f.content {
		n, err = buffer.Write(c)
		r = r + n
		contentRanges[i] = &filepb.ByteRange{
			Start: int32(l),
			End:   int32(r),
		}
		totalContentLength += utf8.RuneCount(c)
		l = r
	}

	var totalImageSize int
	imageRanges := make([]*filepb.ByteRange, len(f.images))
	for i, img := range f.images {
		n, err = buffer.Write(img)
		imageRanges[i] = &filepb.ByteRange{
			Start: int32(l),
			End:   int32(r),
		}
		totalImageSize += n
		l = r
	}

	var totalSectionSize int
	sectionRanges := make([]*filepb.ByteRange, len(f.sections))
	for i, section := range f.sections {
		n, err = section.MarshalToVT()
		sectionRanges[i] = &filepb.ByteRange{
			Start: int32(offset),
			End:   int32(offset + len(section)),
		}
		totalSize += len(section)
		totalSectionSize += n
		l = r
	}

	footer := &filepb.Metadata{
		Version:       1,
		ContentRanges: contentRanges,
	}

	footerSize := footer.SizeVT()

	n = copy(fileBytes[leftOffset:], fileMarker)
	if n != fileMarkerSize {
		return nil, fmt.Errorf("%w: expected %d: got %d", ErrorUnexpectedTailFileMarkerSize, fileMarkerSize, n)
	}

	return fileBytes[:rightOffset], nil
}

func (f *blockFileV1) DecodeContent(index int) (string, error) {
	if !utf8.Valid(f.content[index]) {
		return "", ErrorContentIsNotValidUTF8
	}
	return string(f.content[index]), nil
}

func (f *blockFileV1) DecodeImage(index int) ([]byte, error) {
	imageSize := len(f.images[index])
	b := make([]byte, imageSize)
	n, err := base64.StdEncoding.Decode(b, f.images[index])
	if err != nil {
		return nil, fmt.Errorf("%w:%v", ErrorImageIsNotValidUTF8, err)
	}
	if n != imageSize {
		return nil, fmt.Errorf("%w: expected the entire %d bytes to be base64: decoded only %d bytes", ErrorImageIsNotValidUTF8, imageSize, n)
	}

	return f.images[index], nil
}

func (f *blockFileV1) DecodeSection(index int) (*filepb.SectionNode, error) {
	return f.sections[index], nil
}
