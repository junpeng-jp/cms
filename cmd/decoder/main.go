package main

import (
	"github.com/junpeng.ong/blog/internal/encoding/storage"
	"github.com/junpeng.ong/blog/internal/file"
	"google.golang.org/protobuf/proto"
)

func main() {}

func DecodeFile(b []byte) (*file.File, error) {
	proto.Marshal()
	return storage.DecodeFile(b)
}
