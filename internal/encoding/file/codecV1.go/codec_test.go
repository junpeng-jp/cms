package codecV1

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/junpeng.ong/blog/internal/filepb"
	"github.com/stretchr/testify/assert"
)

type seekableWriter struct {
	buf []byte
	pos int
}

func (m *seekableWriter) Write(p []byte) (n int, err error) {
	if len(p) > cap(m.buf)-m.pos {
		m.buf = append(m.buf, make([]byte, len(p))...)[:m.pos]
	}
	m.buf = m.buf[:m.pos+len(p)]
	n = copy(m.buf[m.pos:], p)
	m.pos += n
	return n, nil
}

func (m *seekableWriter) Seek(offset int64, whence int) (int64, error) {
	newPos, offs := 0, int(offset)
	switch whence {
	case io.SeekStart:
		newPos = offs
	case io.SeekCurrent:
		newPos = m.pos + offs
	case io.SeekEnd:
		newPos = len(m.buf) + offs
	}
	if newPos < 0 {
		return 0, errors.New("negative result pos")
	}
	m.pos = newPos
	return int64(newPos), nil
}

func TestCodecV1Idempotence(t *testing.T) {
	canonicalSection := &filepb.SectionNode{
		Kind: &filepb.SectionNode_HorizontalLayout{
			HorizontalLayout: &filepb.HorizontalLayout{
				BlockContainers: []*filepb.BlockContainer{
					{
						Blocks: []*filepb.BlockNode{
							{
								Kind: &filepb.BlockNode_QuoteBlock{
									QuoteBlock: &filepb.QuoteBlock{
										Block: []*filepb.BlockNode{
											{
												Kind: &filepb.BlockNode_ParagraphBlock{
													ParagraphBlock: &filepb.ParagraphBlock{
														Inline: []*filepb.InlineNode{
															{
																Kind: &filepb.InlineNode_Bold{
																	Bold: &filepb.Bold{
																		Inline: []*filepb.InlineNode{
																			{
																				Kind: &filepb.InlineNode_Underline{
																					Underline: &filepb.Underline{
																						Inline: []*filepb.InlineNode{
																							{
																								Kind: &filepb.InlineNode_Text{
																									Text: &filepb.Text{
																										Start: 0,
																										End:   10,
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																			{
																				Kind: &filepb.InlineNode_Underline{
																					Underline: &filepb.Underline{
																						Inline: []*filepb.InlineNode{
																							{
																								Kind: &filepb.InlineNode_Text{
																									Text: &filepb.Text{
																										Start: 10,
																										End:   20,
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																			{
																				Kind: &filepb.InlineNode_Text{
																					Text: &filepb.Text{
																						Start: 20,
																						End:   30,
																					},
																				},
																			},
																			{
																				Kind: &filepb.InlineNode_Image{
																					Image: &filepb.Image{
																						Start: 30,
																						End:   46,
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	canonicalContent := []byte("1123456789" + "2123456789" + "3123456789" + "fasd3472-_bbf===")
	canonicalMetadata := &filepb.Metadata{
		Version: 1,
		Size:    int64(len(canonicalContent) + canonicalSection.SizeVT()),
		ContentMetadata: &filepb.ByteRange{
			Start: 0,
			End:   int32(len(canonicalContent)),
		},
		SectionMetadata: &filepb.SectionMetadata{
			Ranges: []*filepb.ByteRange{
				{
					Start: int32(len(canonicalContent)),
					End:   int32(len(canonicalContent) + canonicalSection.SizeVT()),
				},
			},
		},
		FileMetadata: &filepb.FileMetadata{
			Name:      "some name",
			CreatedAt: 0,
		},
	}

	testCases := []struct {
		name     string
		offset   int
		section  *filepb.SectionNode
		content  []byte
		metadata *filepb.Metadata
	}{
		{
			name:     "no error: section encoded at offset 0",
			offset:   0,
			section:  canonicalSection,
			content:  canonicalContent,
			metadata: canonicalMetadata,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			writer := &seekableWriter{buf: make([]byte, 64)}

			var err error

			encoder := NewBlockFileEncoderV1(writer, tc.offset)
			err = encoder.EncodeSectionContent(tc.section, tc.content)
			assert.NoError(t, err)
			_, err = encoder.Finalize()
			assert.NoError(t, err)

			metadata := &filepb.Metadata{
				Version:         1,
				Size:            int64(len(writer.buf)),
				ContentMetadata: encoder.GetFinalContentMetadata(),
				SectionMetadata: encoder.GetFinalSectionMetadata(),
				FileMetadata: &filepb.FileMetadata{
					Name:      "some name",
					CreatedAt: 0,
				},
			}

			assert.Equal(t, tc.metadata, metadata)

			reader := bytes.NewReader(writer.buf)

			decoder := NewBlockFileDecoderV1(reader, tc.metadata)
			decoded, err := decoder.DecodeSection(tc.offset)
			assert.NoError(t, err)

			assert.Equal(t, tc.section, decoded)
		})
	}
}
