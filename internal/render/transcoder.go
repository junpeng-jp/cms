package render

import (
	"github.com/junpeng.ong/blog/internal/encoding"
	"github.com/junpeng.ong/blog/internal/filepb"
)

type Transcoder interface {
	ConvertLayoutNodeBlockContainer(*filepb.BlockContainer, int) error
	ConvertLayoutNodeHorizontalLayout(*filepb.HorizontalLayout, int) error
	ConvertLayoutNodeColumnLayout1(*filepb.ColumnLayout1, int) error
	ConvertLayoutNodeColumnLayout2(*filepb.ColumnLayout2, int) error
	ConvertLayoutNodeColumnLayout3(*filepb.ColumnLayout3, int) error
	ConvertLayoutNodeColumnLayout4(*filepb.ColumnLayout4, int) error
	ConvertBlockNodeParagraphBlock(*filepb.ParagraphBlock, int) error
	ConvertBlockNodeDividerBlock(*filepb.DividerBlock, int) error
	ConvertBlockNodeCodeBlock(*filepb.CodeBlock, int) error
	ConvertBlockNodeListBlock(*filepb.ListBlock, int) error
	ConvertBlockNodeTodoListBlock(*filepb.TodoListBlock, int) error
	ConvertBlockNodeQuoteBlock(*filepb.QuoteBlock, int) error
	ConvertInlineNodeText(*filepb.Text, int) error
	ConvertInlineNodeImage(*filepb.Image, int) error
	ConvertInlineNodeLink(*filepb.Link, int) error
	ConvertInlineNodeMention(*filepb.Mention, int) error
	ConvertInlineNodeEquation(*filepb.Equation, int) error
	ConvertInlineNodeCode(*filepb.Code, int) error
	ConvertInlineNodeCustomFormat(*filepb.CustomFormat, int) error
	ConvertInlineNodeBold(*filepb.Bold, int) error
	ConvertInlineNodeItalic(*filepb.Italic, int) error
	ConvertInlineNodeUnderline(*filepb.Underline, int) error
	ConvertInlineNodeStrikethrough(*filepb.Strikethrough, int) error
}

func TranscodeFile(decoder encoding.BlockFileLazyDecoder, transcoder Transcoder) error {
	for i := range decoder.Length() {
		section, err := decoder.DecodeSection(i)
		if err != nil {
			return err
		}
		WalkSection(section, transcoder, 0)
	}
	return nil
}

func WalkSection(section *filepb.SectionNode, transcoder Transcoder, depth int) error {
	var err error
	for _, blocks := range section.Children {
		switch n := blocks.Kind.(type) {
		case *filepb.LayoutAndBlockNode_BlockContainers:
			err = transcoder.ConvertLayoutNodeBlockContainer(n.BlockContainers, depth)
			if err != nil {
				return err
			}
		case *filepb.LayoutAndBlockNode_HorizontalLayout:
			err = transcoder.ConvertLayoutNodeHorizontalLayout(n.HorizontalLayout, depth)
			if err != nil {
				return err
			}
		case *filepb.LayoutAndBlockNode_ColumnLayout_1:
			err = transcoder.ConvertLayoutNodeColumnLayout1(n.ColumnLayout_1, depth)
			if err != nil {
				return err
			}
		case *filepb.LayoutAndBlockNode_ColumnLayout_2:
			err = transcoder.ConvertLayoutNodeColumnLayout2(n.ColumnLayout_2, depth)
			if err != nil {
				return err
			}
		case *filepb.LayoutAndBlockNode_ColumnLayout_3:
			err = transcoder.ConvertLayoutNodeColumnLayout3(n.ColumnLayout_3, depth)
			if err != nil {
				return err
			}
		case *filepb.LayoutAndBlockNode_ColumnLayout_4:
			err = transcoder.ConvertLayoutNodeColumnLayout4(n.ColumnLayout_4, depth)
			if err != nil {
				return err
			}
		case *filepb.LayoutAndBlockNode_ParagraphBlock:
			err = transcoder.ConvertBlockNodeParagraphBlock(n.ParagraphBlock, depth)
			if err != nil {
				return err
			}
		case *filepb.LayoutAndBlockNode_DividerBlock:
			err = transcoder.ConvertBlockNodeDividerBlock(n.DividerBlock, depth)
			if err != nil {
				return err
			}
		case *filepb.LayoutAndBlockNode_CodeBlock:
			err = transcoder.ConvertBlockNodeCodeBlock(n.CodeBlock, depth)
			if err != nil {
				return err
			}
		case *filepb.LayoutAndBlockNode_ListBlock:
			err = transcoder.ConvertBlockNodeListBlock(n.ListBlock, depth)
			if err != nil {
				return err
			}
		case *filepb.LayoutAndBlockNode_TodoListBlock:
			err = transcoder.ConvertBlockNodeTodoListBlock(n.TodoListBlock, depth)
			if err != nil {
				return err
			}
		case *filepb.LayoutAndBlockNode_QuoteBlock:
			err = transcoder.ConvertBlockNodeQuoteBlock(n.QuoteBlock, depth)
			if err != nil {
				return err
			}
		default:
		}
	}
	return nil
}

func WalkBlockNode(block *filepb.BlockNode, transcoder Transcoder, depth int) error {
	var err error
	switch n := block.Kind.(type) {
	case *filepb.BlockNode_ParagraphBlock:
		err = transcoder.ConvertBlockNodeParagraphBlock(n.ParagraphBlock, depth)
		if err != nil {
			return err
		}
	case *filepb.BlockNode_DividerBlock:
		err = transcoder.ConvertBlockNodeDividerBlock(n.DividerBlock, depth)
		if err != nil {
			return err
		}
	case *filepb.BlockNode_CodeBlock:
		err = transcoder.ConvertBlockNodeCodeBlock(n.CodeBlock, depth)
		if err != nil {
			return err
		}
	case *filepb.BlockNode_ListBlock:
		err = transcoder.ConvertBlockNodeListBlock(n.ListBlock, depth)
		if err != nil {
			return err
		}
	case *filepb.BlockNode_TodoListBlock:
		err = transcoder.ConvertBlockNodeTodoListBlock(n.TodoListBlock, depth)
		if err != nil {
			return err
		}
	case *filepb.BlockNode_QuoteBlock:
		err = transcoder.ConvertBlockNodeQuoteBlock(n.QuoteBlock, depth)
		if err != nil {
			return err
		}
	default:
	}
	return nil
}

func WalkInlineNode(inline *filepb.InlineNode, transcoder Transcoder, depth int) error {
	var err error
	switch n := inline.Kind.(type) {
	case *filepb.InlineNode_Text:
		err = transcoder.ConvertInlineNodeText(n.Text, depth)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Image:
		err = transcoder.ConvertInlineNodeImage(n.Image, depth)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Link:
		err = transcoder.ConvertInlineNodeLink(n.Link, depth)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Mention:
		err = transcoder.ConvertInlineNodeMention(n.Mention, depth)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Equation:
		err = transcoder.ConvertInlineNodeEquation(n.Equation, depth)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Code:
		err = transcoder.ConvertInlineNodeCode(n.Code, depth)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_CustomFormat:
		err = transcoder.ConvertInlineNodeCustomFormat(n.CustomFormat, depth)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Bold:
		err = transcoder.ConvertInlineNodeBold(n.Bold, depth)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Italic:
		err = transcoder.ConvertInlineNodeItalic(n.Italic, depth)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Underline:
		err = transcoder.ConvertInlineNodeUnderline(n.Underline, depth)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Strikethrough:
		err = transcoder.ConvertInlineNodeStrikethrough(n.Strikethrough, depth)
		if err != nil {
			return err
		}
	default:
	}
	return nil
}
