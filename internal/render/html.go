package render

import (
	"io"

	"github.com/junpeng.ong/blog/internal/encoding"
	"github.com/junpeng.ong/blog/internal/filepb"
)

const (
	tab = "  "
)

type HtmlTranscoder struct {
	decoder encoding.BlockFileLazyDecoder
	writer  io.Writer
	pos     int
}

func NewHtmlTranscoder(writer io.Writer) *HtmlTranscoder {
	return &HtmlTranscoder{writer: writer}
}

func (t *HtmlTranscoder) ConvertSectionNodeBlockContainer(node *filepb.BlockContainer, depth int) (err error) {
	// n, err := writeOpeningTag(t.writer, tagDiv, depth, createClassAnnotation("col"))
	// t.pos += n
	// if err != nil {
	// 	return err
	// }
	var n int
	for _, block := range node.Blocks {
		n, err = writeOpeningTag(t.writer, tagDiv, depth+1, createClassAnnotation("block"))
		t.pos += n
		if err != nil {
			return err
		}
		err := WalkBlockNode(block, t, depth+2)
		if err != nil {
			return err
		}
		n, err = writeClosingTag(t.writer, tagDiv, depth+1)
		t.pos += n
		if err != nil {
			return err
		}
	}
	// n, err = writeClosingTag(t.writer, tagDiv, depth)
	// t.pos += n
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (t *HtmlTranscoder) ConvertSectionNodeHorizontalLayout(node *filepb.HorizontalLayout, depth int) (err error) {
	n, err := writeOpeningTag(t.writer, tagDiv, depth, createClassAnnotation("row"))
	t.pos += n
	if err != nil {
		return err
	}
	for _, blockContainer := range node.GetBlockContainers() {
		err := t.ConvertSectionNodeBlockContainer(blockContainer, depth+1)
		if err != nil {
			return nil
		}
	}
	n, err = writeClosingTag(t.writer, tagDiv, depth)
	t.pos += n
	if err != nil {
		return err
	}
	return nil
}

func (t *HtmlTranscoder) ConvertSectionNodeColumnLayout1(node *filepb.ColumnLayout1, depth int) (err error) {
	n, err := writeOpeningTag(t.writer, tagDiv, depth, createClassAnnotation("col"))
	t.pos += n
	if err != nil {
		return err
	}
	err = t.ConvertSectionNodeBlockContainer(node.GetBlockContainer(), depth+1)
	if err != nil {
		return nil
	}
	n, err = writeClosingTag(t.writer, tagDiv, depth)
	t.pos += n
	if err != nil {
		return err
	}
	return nil
}

func (t *HtmlTranscoder) ConvertSectionNodeColumnLayout2(node *filepb.ColumnLayout2, depth int) (err error) {
	n, err := writeOpeningTag(t.writer, tagDiv, depth, createClassAnnotation("col"))
	t.pos += n
	if err != nil {
		return err
	}
	for _, blockContainer := range node.GetBlockContainers() {
		err := t.ConvertSectionNodeBlockContainer(blockContainer, depth+1)
		if err != nil {
			return nil
		}
	}
	n, err = writeClosingTag(t.writer, tagDiv, depth)
	t.pos += n
	if err != nil {
		return err
	}
	return nil
}

func (t *HtmlTranscoder) ConvertSectionNodeColumnLayout3(node *filepb.ColumnLayout3, depth int) (err error) {
	n, err := writeOpeningTag(t.writer, tagDiv, depth, createClassAnnotation("col"))
	t.pos += n
	if err != nil {
		return err
	}
	for _, blockContainer := range node.GetBlockContainers() {
		err := t.ConvertSectionNodeBlockContainer(blockContainer, depth+1)
		if err != nil {
			return nil
		}
	}
	n, err = writeClosingTag(t.writer, tagDiv, depth)
	t.pos += n
	if err != nil {
		return err
	}
	return nil
}

func (t *HtmlTranscoder) ConvertSectionNodeColumnLayout4(node *filepb.ColumnLayout4, depth int) (err error) {
	n, err := writeOpeningTag(t.writer, tagDiv, depth, createClassAnnotation("col"))
	t.pos += n
	if err != nil {
		return err
	}
	for _, blockContainer := range node.GetBlockContainers() {
		err := t.ConvertSectionNodeBlockContainer(blockContainer, depth+1)
		if err != nil {
			return nil
		}
	}
	n, err = writeClosingTag(t.writer, tagDiv, depth)
	t.pos += n
	if err != nil {
		return err
	}
	return nil
}

func (t *HtmlTranscoder) ConvertBlockNodeParagraphBlock(node *filepb.ParagraphBlock, depth int) (err error) {
	n, err := writeOpeningTag(t.writer, tagParagraph, depth)
	t.pos += n
	if err != nil {
		return err
	}
	for _, inlineNode := range node.GetInline() {
		err := WalkInlineNode(inlineNode, t, depth+1)
		if err != nil {
			return nil
		}
	}
	n, err = writeClosingTag(t.writer, tagParagraph, depth)
	t.pos += n
	if err != nil {
		return err
	}
	return nil
}

func (t *HtmlTranscoder) ConvertBlockNodeDividerBlock(node *filepb.DividerBlock, depth int) (err error) {
	t.writer.Write([]byte("<divider>"))
	t.writer.Write([]byte("</divider>"))
	return nil
}

func (t *HtmlTranscoder) ConvertBlockNodeCodeBlock(node *filepb.CodeBlock, depth int) (err error) {
	t.writer.Write([]byte("<code block)>"))
	t.writer.Write([]byte("</code block)>"))
	return nil
}

func (t *HtmlTranscoder) ConvertBlockNodeListBlock(node *filepb.ListBlock, depth int) (err error) {
	t.writer.Write([]byte("<ol>"))
	t.writer.Write([]byte("</ol>"))
	return nil
}

func (t *HtmlTranscoder) ConvertBlockNodeTodoListBlock(node *filepb.TodoListBlock, depth int) (err error) {
	t.writer.Write([]byte("<todo>"))
	t.writer.Write([]byte("</todo>"))
	return nil
}

func (t *HtmlTranscoder) ConvertBlockNodeQuoteBlock(node *filepb.QuoteBlock, depth int) (err error) {
	n, err := writeOpeningTag(t.writer, tagBlockquote, depth)
	t.pos += n
	if err != nil {
		return err
	}
	for _, blockNode := range node.GetBlock() {
		err := WalkBlockNode(blockNode, t, depth+1)
		if err != nil {
			return nil
		}
	}
	n, err = writeClosingTag(t.writer, tagBlockquote, depth)
	t.pos += n
	if err != nil {
		return err
	}
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeText(node *filepb.Text, depth int) (err error) {
	content, err := t.decoder.DecodeContent(int(node.Start), int(node.End-node.Start))
	if err != nil {
		return err
	}
	n, err := writeContent(t.writer, content, depth)
	t.pos += n
	if err != nil {
		return err
	}
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeImage(node *filepb.Image, depth int) (err error) {
	n, err := writeOpeningTag(t.writer, tagImage, depth)
	t.pos += n
	if err != nil {
		return err
	}
	img, err := t.decoder.DecodeBase64Image(int(node.Start), int(node.End-node.Start))
	if err != nil {
		return err
	}
	n, err = writeContent(t.writer, img, depth+1)
	t.pos += n
	if err != nil {
		return err
	}
	n, err = writeClosingTag(t.writer, tagImage, depth)
	t.pos += n
	if err != nil {
		return err
	}
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeLink(node *filepb.Link, depth int) (err error) {
	t.writer.Write([]byte("<a>"))
	t.writer.Write([]byte("</a>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeMention(node *filepb.Mention, depth int) (err error) {
	t.writer.Write([]byte("<mention>"))
	t.writer.Write([]byte("</mention>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeEquation(node *filepb.Equation, depth int) (err error) {
	t.writer.Write([]byte("<equation>"))
	t.writer.Write([]byte("</equation>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeCode(node *filepb.Code, depth int) (err error) {
	t.writer.Write([]byte("<code>"))
	t.writer.Write([]byte("</code>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeCustomFormat(node *filepb.CustomFormat, depth int) (err error) {
	t.writer.Write([]byte("<custom>"))
	t.writer.Write([]byte("</custom>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeBold(node *filepb.Bold, depth int) (err error) {
	n, err := writeOpeningTag(t.writer, tagUnderline, depth)
	t.pos += n
	if err != nil {
		return err
	}
	for _, inlineNode := range node.GetInline() {
		err := WalkInlineNode(inlineNode, t, depth+1)
		if err != nil {
			return nil
		}
	}
	n, err = writeClosingTag(t.writer, tagUnderline, depth)
	t.pos += n
	if err != nil {
		return err
	}
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeItalic(node *filepb.Italic, depth int) (err error) {
	n, err := writeOpeningTag(t.writer, tagItalic, depth)
	t.pos += n
	if err != nil {
		return err
	}
	for _, inlineNode := range node.GetInline() {
		err := WalkInlineNode(inlineNode, t, depth+1)
		if err != nil {
			return nil
		}
	}
	n, err = writeClosingTag(t.writer, tagItalic, depth)
	t.pos += n
	if err != nil {
		return err
	}
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeUnderline(node *filepb.Underline, depth int) (err error) {
	n, err := writeOpeningTag(t.writer, tagUnderline, depth)
	t.pos += n
	if err != nil {
		return err
	}
	for _, inlineNode := range node.GetInline() {
		err := WalkInlineNode(inlineNode, t, depth+1)
		if err != nil {
			return nil
		}
	}
	n, err = writeClosingTag(t.writer, tagUnderline, depth)
	t.pos += n
	if err != nil {
		return err
	}
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeStrikethrough(node *filepb.Strikethrough, depth int) (err error) {
	t.writer.Write([]byte("<s>"))
	t.writer.Write([]byte("</s>"))
	return nil
}
