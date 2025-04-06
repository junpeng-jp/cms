package render

import (
	"testing"

	"github.com/junpeng.ong/blog/internal/filepb"
	"github.com/junpeng.ong/blog/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestHtmlTranscoder(t *testing.T) {
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

	bout := make([]byte, 64)
	htmlTranscoder := HtmlTranscoder{
		writer: &testutils.SeekableWriter{Buf: bout},
	}
	err := WalkSection(canonicalSection, &htmlTranscoder)
	assert.NoError(t, err)

	t.Logf("%s", bout)
}
