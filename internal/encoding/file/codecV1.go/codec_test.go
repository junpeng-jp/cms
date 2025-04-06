package codecV1

import (
	"bytes"
	"testing"

	"github.com/junpeng.ong/blog/internal/encoding/common"
	"github.com/junpeng.ong/blog/internal/encoding/utils"
	"github.com/junpeng.ong/blog/internal/filepb"
	"github.com/stretchr/testify/assert"
)

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
		ContentMetadata: &filepb.ByteRange{
			Start: common.FileMarkerSize,
			End:   common.FileMarkerSize + int32(len(canonicalContent)),
		},
		SectionMetadata: &filepb.SectionMetadata{
			Ranges: []*filepb.ByteRange{
				{
					Start: common.FileMarkerSize + int32(len(canonicalContent)),
					End:   common.FileMarkerSize + int32(len(canonicalContent)+canonicalSection.SizeVT()),
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
			writer := utils.NewSeekableWriter(make([]byte, 64))

			var err error

			encoder := NewBlockFileEncoderV1(writer)
			err = encoder.Init()
			assert.NoError(t, err)
			err = encoder.EncodeSectionContent(tc.section, tc.content)
			assert.NoError(t, err)
			_, err = encoder.Finalize("some name")
			assert.NoError(t, err)

			t.Logf("%s", writer.GetBuffer())
			t.Logf("size: %d", len(writer.GetBuffer()))

			reader := bytes.NewReader(writer.GetBuffer())

			err = common.VerifyFileMarkers(reader)
			assert.NoError(t, err)
			_, metadata, err := common.DecodeMetadata(reader)
			assert.NoError(t, err)
			metadata.FileMetadata.CreatedAt = 0
			assert.Equal(t, tc.metadata, metadata)

			decoder := NewBlockFileDecoderV1(reader, tc.metadata)
			decoded, err := decoder.DecodeSection(tc.offset)
			assert.NoError(t, err)

			assert.Equal(t, tc.section, decoded)
		})
	}
}
