package main

import (
	"bytes"
	"log"
	"os"

	"github.com/junpeng.ong/blog/internal/encoding"
	"github.com/junpeng.ong/blog/internal/filepb"
)

var exampleContent = []byte("1123456789" + "2123456789" + "3123456789" + "fasd3472-_bbf===")
var exampleFile = &filepb.SectionNode{
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
																									Start: 1,
																									End:   262144,
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
																									Start: 1,
																									End:   262144,
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
																					Start: 1,
																					End:   262144,
																				},
																			},
																		},
																		{
																			Kind: &filepb.InlineNode_Text{
																				Text: &filepb.Text{
																					Start: 1,
																					End:   262144,
																				},
																			},
																		},
																		{
																			Kind: &filepb.InlineNode_Text{
																				Text: &filepb.Text{
																					Start: 1,
																					End:   262144,
																				},
																			},
																		}, {
																			Kind: &filepb.InlineNode_Text{
																				Text: &filepb.Text{
																					Start: 1,
																					End:   262144,
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

func main() {
	buffer := bytes.NewBuffer(make([]byte, 64))
	encoder, err := encoding.NewBlockFileEncoder(buffer)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = encoder.Init()
	if err != nil {
		log.Fatalf("failed to encode file: %s\n", err)
		return
	}
	err = encoder.EncodeSectionContent(exampleFile, exampleContent)
	if err != nil {
		log.Fatalf("failed to encode file: %s\n", err)
		return
	}
	_, err = encoder.Finalize("example")
	if err != nil {
		log.Fatalf("failed to encode file: %s\n", err)
		return
	}
	if err := os.WriteFile("build/data/example.bloc", buffer.Bytes(), os.ModePerm); err != nil {
		log.Fatalf("failed to write file: %s\n", err)
	}
}
