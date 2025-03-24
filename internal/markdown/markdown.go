package markdown

import (
	"io"

	"github.com/junpeng.ong/blog/internal/pb/block"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
)

type Parser interface {
	Parse(source []byte)
}

type Renderer interface {
	Render(w io.Writer, source []byte, blocks []block.Block)
}

type MarkdownConverter struct {
	md goldmark.Markdown
}

func (c *MarkdownConverter) Parse(source []byte) (*block.Block, error) {
	// parse markdown into a ast format

	// convert ast to proto struct

	// marshal proto struct to binary
	//
	switch source.Kind() {
	case ast.KindDocument:
	}

	return nil, nil
}

func (c *MarkdownConverter) Render(w io.Writer, source []byte, blocks []block.Block) (string, error) {
	//
	switch source.Kind() {
	case ast.KindDocument:
	}

	return "", nil
}
