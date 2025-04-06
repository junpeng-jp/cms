package render

import (
	"github.com/junpeng.ong/blog/internal/encoding"
	"github.com/junpeng.ong/blog/internal/filepb"
)

type Transcoder interface {
	ConvertSectionNodeBlockContainer(*filepb.BlockContainer) error
	ConvertSectionNodeHorizontalLayout(*filepb.HorizontalLayout) error
	ConvertSectionNodeColumnLayout1(*filepb.ColumnLayout1) error
	ConvertSectionNodeColumnLayout2(*filepb.ColumnLayout2) error
	ConvertSectionNodeColumnLayout3(*filepb.ColumnLayout3) error
	ConvertSectionNodeColumnLayout4(*filepb.ColumnLayout4) error
	ConvertBlockNodeParagraphBlock(*filepb.ParagraphBlock) error
	ConvertBlockNodeDividerBlock(*filepb.DividerBlock) error
	ConvertBlockNodeCodeBlock(*filepb.CodeBlock) error
	ConvertBlockNodeListBlock(*filepb.ListBlock) error
	ConvertBlockNodeTodoListBlock(*filepb.TodoListBlock) error
	ConvertBlockNodeQuoteBlock(*filepb.QuoteBlock) error
	ConvertInlineNodeText(*filepb.Text) error
	ConvertInlineNodeImage(*filepb.Image) error
	ConvertInlineNodeLink(*filepb.Link) error
	ConvertInlineNodeMention(*filepb.Mention) error
	ConvertInlineNodeEquation(*filepb.Equation) error
	ConvertInlineNodeCode(*filepb.Code) error
	ConvertInlineNodeCustomFormat(*filepb.CustomFormat) error
	ConvertInlineNodeBold(*filepb.Bold) error
	ConvertInlineNodeItalic(*filepb.Italic) error
	ConvertInlineNodeUnderline(*filepb.Underline) error
	ConvertInlineNodeStrikethrough(*filepb.Strikethrough) error
}

func TranscodeFile(decoder encoding.BlockFileLazyDecoder, transcoder Transcoder) error {
	for i := range decoder.Length() {
		section, err := decoder.DecodeSection(i)
		if err != nil {
			return err
		}
		WalkSection(section, transcoder)
	}
	return nil
}

func WalkSection(section *filepb.SectionNode, transcoder Transcoder) error {
	var err error
	switch n := section.Kind.(type) {
	case *filepb.SectionNode_BlockContainers:
		err = transcoder.ConvertSectionNodeBlockContainer(n.BlockContainers)
		if err != nil {
			return err
		}
	case *filepb.SectionNode_HorizontalLayout:
		err = transcoder.ConvertSectionNodeHorizontalLayout(n.HorizontalLayout)
		if err != nil {
			return err
		}

	case *filepb.SectionNode_ColumnLayout_1:
		err = transcoder.ConvertSectionNodeColumnLayout1(n.ColumnLayout_1)
		if err != nil {
			return err
		}
	case *filepb.SectionNode_ColumnLayout_2:
		err = transcoder.ConvertSectionNodeColumnLayout2(n.ColumnLayout_2)
		if err != nil {
			return err
		}
	case *filepb.SectionNode_ColumnLayout_3:
		err = transcoder.ConvertSectionNodeColumnLayout3(n.ColumnLayout_3)
		if err != nil {
			return err
		}
	case *filepb.SectionNode_ColumnLayout_4:
		err = transcoder.ConvertSectionNodeColumnLayout4(n.ColumnLayout_4)
		if err != nil {
			return err
		}
	default:

	}
	return nil
}

func WalkBlockNode(block *filepb.BlockNode, transcoder Transcoder) error {
	var err error
	switch n := block.Kind.(type) {
	case *filepb.BlockNode_ParagraphBlock:
		err = transcoder.ConvertBlockNodeParagraphBlock(n.ParagraphBlock)
		if err != nil {
			return err
		}
	case *filepb.BlockNode_DividerBlock:
		err = transcoder.ConvertBlockNodeDividerBlock(n.DividerBlock)
		if err != nil {
			return err
		}
	case *filepb.BlockNode_CodeBlock:
		err = transcoder.ConvertBlockNodeCodeBlock(n.CodeBlock)
		if err != nil {
			return err
		}
	case *filepb.BlockNode_ListBlock:
		err = transcoder.ConvertBlockNodeListBlock(n.ListBlock)
		if err != nil {
			return err
		}
	case *filepb.BlockNode_TodoListBlock:
		err = transcoder.ConvertBlockNodeTodoListBlock(n.TodoListBlock)
		if err != nil {
			return err
		}
	case *filepb.BlockNode_QuoteBlock:
		err = transcoder.ConvertBlockNodeQuoteBlock(n.QuoteBlock)
		if err != nil {
			return err
		}
	}
	return nil
}

func WalkInlineNode(inline *filepb.InlineNode, transcoder Transcoder) error {
	var err error
	switch n := inline.Kind.(type) {
	case *filepb.InlineNode_Text:
		err = transcoder.ConvertInlineNodeText(n.Text)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Image:
		err = transcoder.ConvertInlineNodeImage(n.Image)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Link:
		err = transcoder.ConvertInlineNodeLink(n.Link)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Mention:
		err = transcoder.ConvertInlineNodeMention(n.Mention)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Equation:
		err = transcoder.ConvertInlineNodeEquation(n.Equation)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Code:
		err = transcoder.ConvertInlineNodeCode(n.Code)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_CustomFormat:
		err = transcoder.ConvertInlineNodeCustomFormat(n.CustomFormat)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Bold:
		err = transcoder.ConvertInlineNodeBold(n.Bold)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Italic:
		err = transcoder.ConvertInlineNodeItalic(n.Italic)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Underline:
		err = transcoder.ConvertInlineNodeUnderline(n.Underline)
		if err != nil {
			return err
		}
	case *filepb.InlineNode_Strikethrough:
		err = transcoder.ConvertInlineNodeStrikethrough(n.Strikethrough)
		if err != nil {
			return err
		}
	}
	return nil
}
