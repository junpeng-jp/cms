package file

import (
	"github.com/junpeng.ong/blog/internal/pb/contentpb"
	"github.com/junpeng.ong/blog/internal/pb/filepb"
)

type File struct {
	Content *contentpb.Content
	Footer  *filepb.Footer
}

func NewFile() *File {}

func NewFileFromBinary() *File {}

func (f *File) ToBinary() ([]byte, error) {
	f.Content.SizeVT()
	b, err := f.Content.MarshalVT()
	if err != nil {
		return nil, err
	}

	return b, nil
}
