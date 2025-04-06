package render

import (
	"io"

	"github.com/junpeng.ong/blog/internal/filepb"
)

type HtmlTranscoder struct {
	writer io.Writer
}

func NewHtmlTranscoder(writer io.Writer) *HtmlTranscoder {
	return &HtmlTranscoder{writer: writer}
}

func (t *HtmlTranscoder) ConvertSectionNodeBlockContainer(node *filepb.BlockContainer) error {
	t.writer.Write([]byte("<container>\n"))
	for _, block := range node.Blocks {
		t.writer.Write([]byte("<div>"))
		err := WalkBlockNode(block, t)
		if err != nil {
			return err
		}
		t.writer.Write([]byte("</div>"))
	}
	t.writer.Write([]byte("</container>"))
	return nil
}

func (t *HtmlTranscoder) ConvertSectionNodeHorizontalLayout(node *filepb.HorizontalLayout) error {
	t.writer.Write([]byte("<horizontal>"))
	for _, blockContainer := range node.GetBlockContainers() {
		err := t.ConvertSectionNodeBlockContainer(blockContainer)
		if err != nil {
			return nil
		}
	}
	t.writer.Write([]byte("</horizontal>"))
	return nil
}

func (t *HtmlTranscoder) ConvertSectionNodeColumnLayout1(node *filepb.ColumnLayout1) error {
	t.writer.Write([]byte("<column1>"))
	t.writer.Write([]byte("</column1>"))
	return nil
}

func (t *HtmlTranscoder) ConvertSectionNodeColumnLayout2(node *filepb.ColumnLayout2) error {
	t.writer.Write([]byte("<column2>"))
	t.writer.Write([]byte("</column2>"))
	return nil
}

func (t *HtmlTranscoder) ConvertSectionNodeColumnLayout3(node *filepb.ColumnLayout3) error {
	t.writer.Write([]byte("<column3>"))
	t.writer.Write([]byte("</column3>"))
	return nil
}

func (t *HtmlTranscoder) ConvertSectionNodeColumnLayout4(node *filepb.ColumnLayout4) error {
	t.writer.Write([]byte("<column4>"))
	t.writer.Write([]byte("</column4>"))
	return nil
}

func (t *HtmlTranscoder) ConvertBlockNodeParagraphBlock(node *filepb.ParagraphBlock) error {
	t.writer.Write([]byte("<paragraph>"))
	t.writer.Write([]byte("</paragraph>"))
	return nil
}

func (t *HtmlTranscoder) ConvertBlockNodeDividerBlock(node *filepb.DividerBlock) error {
	t.writer.Write([]byte("<divider>"))
	t.writer.Write([]byte("</divider>"))
	return nil
}

func (t *HtmlTranscoder) ConvertBlockNodeCodeBlock(node *filepb.CodeBlock) error {
	t.writer.Write([]byte("<code block)>"))
	t.writer.Write([]byte("</code block)>"))
	return nil
}

func (t *HtmlTranscoder) ConvertBlockNodeListBlock(node *filepb.ListBlock) error {
	t.writer.Write([]byte("<ol>"))
	t.writer.Write([]byte("</ol>"))
	return nil
}

func (t *HtmlTranscoder) ConvertBlockNodeTodoListBlock(node *filepb.TodoListBlock) error {
	t.writer.Write([]byte("<todo>"))
	t.writer.Write([]byte("</todo>"))
	return nil
}

func (t *HtmlTranscoder) ConvertBlockNodeQuoteBlock(node *filepb.QuoteBlock) error {
	t.writer.Write([]byte("<quote>"))
	t.writer.Write([]byte("</quote>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeText(node *filepb.Text) error {
	t.writer.Write([]byte("<>"))
	t.writer.Write([]byte("</>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeImage(node *filepb.Image) error {
	t.writer.Write([]byte("<img>"))
	t.writer.Write([]byte("</img>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeLink(node *filepb.Link) error {
	t.writer.Write([]byte("<a>"))
	t.writer.Write([]byte("</a>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeMention(node *filepb.Mention) error {
	t.writer.Write([]byte("<mention>"))
	t.writer.Write([]byte("</mention>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeEquation(node *filepb.Equation) error {
	t.writer.Write([]byte("<equation>"))
	t.writer.Write([]byte("</equation>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeCode(node *filepb.Code) error {
	t.writer.Write([]byte("<code>"))
	t.writer.Write([]byte("</code>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeCustomFormat(node *filepb.CustomFormat) error {
	t.writer.Write([]byte("<custom>"))
	t.writer.Write([]byte("</custom>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeBold(node *filepb.Bold) error {
	t.writer.Write([]byte("<b>"))
	t.writer.Write([]byte("</b>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeItalic(node *filepb.Italic) error {
	t.writer.Write([]byte("<i>"))
	t.writer.Write([]byte("</i>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeUnderline(node *filepb.Underline) error {
	t.writer.Write([]byte("<u>"))
	t.writer.Write([]byte("</u>"))
	return nil
}

func (t *HtmlTranscoder) ConvertInlineNodeStrikethrough(node *filepb.Strikethrough) error {
	t.writer.Write([]byte("<s>"))
	t.writer.Write([]byte("</s>"))
	return nil
}
