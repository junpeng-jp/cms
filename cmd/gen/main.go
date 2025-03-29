package main

import (
	"log"
	"os"

	"github.com/junpeng.ong/blog/internal/encoding/storage"
	"github.com/junpeng.ong/blog/internal/file"
	"github.com/junpeng.ong/blog/internal/file/contentpb"
)

var exampleFile = &file.File{
	ContentList: []*contentpb.Content{
		{
			Kind: &contentpb.Content_Text{
				Text: &contentpb.Text{
					Fragments: []*contentpb.TextFragment{
						{
							Annotation: "",
							Text:       "",
						},
					},
				},
			},
		},
		{
			Kind: &contentpb.Content_Text{
				Text: &contentpb.Text{
					Fragments: []*contentpb.TextFragment{
						{
							Annotation: "",
							Text:       "a",
						},
						{
							Annotation: "",
							Text:       "a",
						},
						{
							Annotation: "",
							Text:       "a",
						},
						{
							Annotation: "",
							Text:       "a",
						},
					},
				},
			},
		},
		{
			Kind: &contentpb.Content_Text{
				Text: &contentpb.Text{
					Fragments: []*contentpb.TextFragment{
						{
							Annotation: "bold",
							Text:       "some unicode text",
						},
						{
							Annotation: "italic",
							Text:       "some other text",
						},
					},
				},
			},
		},
		{
			Kind: &contentpb.Content_Link{
				Link: &contentpb.Link{
					Url:  "https://google.com",
					Text: "link to url",
				},
			},
		},
		{
			Kind: &contentpb.Content_Mention{
				Mention: &contentpb.Mention{
					UserId: "my_user",
				},
			},
		},
		{
			Kind: &contentpb.Content_Code{
				Code: &contentpb.Code{
					Text: "let a = 10",
				},
			},
		},
		{
			Kind: &contentpb.Content_Equation{
				Equation: &contentpb.Equation{
					Expression: "a + b",
				},
			},
		},
		{
			Kind: &contentpb.Content_Image{
				Image: &contentpb.Image{
					Url: "https://domain.com/myimage",
				},
			},
		},
		{
			Kind: &contentpb.Content_EmbeddedImage{
				EmbeddedImage: &contentpb.EmbeddedImage{
					Base64: "asdfasdfasdf",
				},
			},
		},
	},
}

func main() {
	codec := &storage.StorageCodec{}
	b, err := codec.Encode(exampleFile)
	if err != nil {
		log.Fatalf("failed to encode file: %s\n", err)
		return
	}
	if err := os.WriteFile("build/data/example.bloc", b, os.ModePerm); err != nil {
		log.Fatalf("failed to write file: %s\n", err)
	}
}
