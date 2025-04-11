package codecV1

import (
	"encoding/json"
	"testing"

	"github.com/junpeng.ong/blog/internal/filepb"
)

var global *filepb.SectionNode
var globalJson *map[string]any
var globalProto *filepb.SectionNode

func BenchmarkDecoderJsonParsing(b *testing.B) {
	var jsonBytes = []byte("{\"section\":{\"horizontal\":[{\"block_container\": [{\"quote\":[{\"bold\":{\"underline\":{\"text\":{\"start\":1,\"end\":262144}}}},{\"bold\":{\"underline\":{\"text\":{\"start\":1,\"end\":262144}}}},{\"text\":{\"start\":1,\"end\":262144}},{\"text\":{\"start\":1,\"end\":262144}},{\"text\":{\"start\":1,\"end\":262144}},{\"text\":{\"start\":1,\"end\":262144}}]}]}]}}")
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		var section map[string]any
		_ = json.Unmarshal(jsonBytes, &section)
		globalJson = &section
	}
}

func BenchmarkDecoderProtobufParsing(b *testing.B) {
	var protoBytes = []byte{
		0x8a, 0x08, 0x57, 0x82, 0x08, 0x54, 0x82, 0x01, 0x51, 0x9a, 0x04, 0x4e, 0x82, 0x01, 0x4b, 0x8a,
		0x01, 0x48, 0x0a, 0x46, 0x4a, 0x44, 0x0a, 0x0c, 0x5a, 0x0a, 0x0a, 0x08, 0x12, 0x06, 0x08, 0x01,
		0x10, 0x80, 0x80, 0x10, 0x0a, 0x0c, 0x5a, 0x0a, 0x0a, 0x08, 0x12, 0x06, 0x08, 0x01, 0x10, 0x80,
		0x80, 0x10, 0x0a, 0x08, 0x12, 0x06, 0x08, 0x01, 0x10, 0x80, 0x80, 0x10, 0x0a, 0x08, 0x12, 0x06,
		0x08, 0x01, 0x10, 0x80, 0x80, 0x10, 0x0a, 0x08, 0x12, 0x06, 0x08, 0x01, 0x10, 0x80, 0x80, 0x10,
		0x0a, 0x08, 0x12, 0x06, 0x08, 0x01, 0x10, 0x80, 0x80, 0x10,
	}
	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		var section filepb.SectionNode
		_ = section.UnmarshalVT(protoBytes)
		globalProto = &section
	}
}

func BenchmarkDecoderAllocations(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		global = NewMessage()
	}
}

func NewMessage() *filepb.SectionNode {
	return &filepb.SectionNode{
		Children: []*filepb.LayoutAndBlockNode{
			{
				Kind: &filepb.LayoutAndBlockNode_HorizontalLayout{
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
																						Kind: &filepb.InlineNode_Italic{
																							Italic: &filepb.Italic{
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
	}
}

// func TestPrintBytes(t *testing.T) {
// 	bs, err := NewMessage().MarshalVT()
// 	assert.NoError(t, err)

// 	var builder strings.Builder
// 	for i, b := range bs {
// 		builder.WriteString(fmt.Sprintf("0x%02x", b))
// 		if (i+1)%16 == 0 {
// 			builder.WriteRune(',')
// 			builder.WriteRune('\n')
// 		} else {
// 			builder.WriteRune(',')
// 			builder.WriteRune(' ')
// 		}
// 	}
// 	t.Log(builder.String())
// }
