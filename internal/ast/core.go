package ast

import "github.com/junpeng.ong/blog/internal/pb/block"

type Parser interface {
	Parse([]byte) (block.Block, error)
}

type Renderer interface {
	Render(block.Block) ([]byte, error)
}
