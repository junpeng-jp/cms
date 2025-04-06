package render

import (
	"io"

	"github.com/junpeng.ong/blog/internal/encoding"
	"github.com/junpeng.ong/blog/internal/filepb"
)

type NodeHandler interface {
	HandleBlockContainer(*filepb.BlockContainer) error
	HandleHorizontalLayout(*filepb.HorizontalLayout) error
	HandleColumnLayout1(*filepb.ColumnLayout1) error
	HandleColumnLayout2(*filepb.ColumnLayout2) error
	HandleColumnLayout3(*filepb.ColumnLayout3) error
	HandleColumnLayout4(*filepb.ColumnLayout4) error
}

func TranscodeFileLazy(reader io.ReadSeeker, writer io.WriteSeeker) error {
	decoder, err := encoding.NewBlockFileLazyDecoder(reader)
	if err != nil {
		return err
	}

	for i := range decoder.Length() {
		section, err := decoder.DecodeSection(i)
		if err != nil {
			return err
		}

		WalkSection(section, handler)
	}

	return nil
}

func WalkSection(section *filepb.SectionNode, handler NodeHandler) error {
	var err error
	switch n := section.Kind.(type) {
	case *filepb.SectionNode_BlockContainers:
		err = handler.HandleBlockContainer(n.BlockContainers)
		if err != nil {
			return err
		}
		for i, block := range n.BlockContainers.GetBlocks() {
			err = WalkBlockNode(block, handler)
			if err != nil {
				return err
			}
		}
	case *filepb.SectionNode_HorizontalLayout:
		err = handler.HandleHorizontalLayout(n.HorizontalLayout)
		for i, block_container := range n.HorizontalLayout.GetBlockContainers() {
			err = handler.HandleBlockContainer(block_container)
			if err != nil {
				return err
			}
			for j, block := range block_container.GetBlocks() {
				err = WalkBlockNode(block, handler)
				if err != nil {
					return err
				}
			}
		}

	case *filepb.SectionNode_ColumnLayout_1:
		err = handler.HandleColumnLayout1(n.ColumnLayout_1)
		if err != nil {
			return err
		}
		err = handler.HandleBlockContainer(n.ColumnLayout_1.BlockContainer)
		if err != nil {
			return err
		}
		for i, block := range n.ColumnLayout_1.BlockContainer.GetBlocks() {
			err = WalkBlockNode(block, handler)
			if err != nil {
				return err
			}
		}
	case *filepb.SectionNode_ColumnLayout_2:
		err = handler.HandleColumnLayout2(n.ColumnLayout_2)
		if err != nil {
			return err
		}
		for i, block_container := range n.ColumnLayout_2.GetBlockContainers() {
			err = handler.HandleBlockContainer(block_container)
			if err != nil {
				return err
			}
			for j, block := range block_container.GetBlocks() {
				err = WalkBlockNode(block, handler)
				if err != nil {
					return err
				}
			}
		}
	case *filepb.SectionNode_ColumnLayout_3:
		err = handler.HandleColumnLayout3(n.ColumnLayout_3)
		if err != nil {
			return err
		}
		for i, block_container := range n.ColumnLayout_3.GetBlockContainers() {
			err = handler.HandleBlockContainer(block_container)
			if err != nil {
				return err
			}
			for j, block := range block_container.GetBlocks() {
				err = WalkBlockNode(block, handler)
				if err != nil {
					return err
				}
			}
		}
	case *filepb.SectionNode_ColumnLayout_4:
		err = handler.HandleColumnLayout4(n.ColumnLayout_4)
		if err != nil {
			return err
		}
		for i, block_container := range n.ColumnLayout_4.GetBlockContainers() {
			err = handler.HandleBlockContainer(block_container)
			if err != nil {
				return err
			}
			for j, block := range block_container.GetBlocks() {
				err = WalkBlockNode(block, handler)
				if err != nil {
					return err
				}
			}
		}
	default:

	}
	return nil
}

func WalkBlockNode(block *filepb.BlockNode, handler NodeHandler) error {}

func WalkInlineNode(block *filepb.InlineNode, handler NodeHandler) error {}
