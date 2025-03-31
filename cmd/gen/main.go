package main

import (
	"log"
	"os"

	"github.com/junpeng.ong/blog/internal/encoding/storage"
	"github.com/junpeng.ong/blog/internal/file"
	"github.com/junpeng.ong/blog/internal/file/filepb"
)

var exampleFile = &file.File{
	ContentList: []*filepb.Content{
		{
			Kind: &filepb.Content_Text{
				Text: &filepb.Text{
					Fragments: []*filepb.TextFragment{
						{
							Annotation: "",
							Text:       "",
						},
					},
				},
			},
		},
		{
			Kind: &filepb.Content_Text{
				Text: &filepb.Text{
					Fragments: []*filepb.TextFragment{
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
			Kind: &filepb.Content_Text{
				Text: &filepb.Text{
					Fragments: []*filepb.TextFragment{
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
			Kind: &filepb.Content_Link{
				Link: &filepb.Link{
					Url:  "https://google.com",
					Text: "link to url",
				},
			},
		},
		{
			Kind: &filepb.Content_Mention{
				Mention: &filepb.Mention{
					UserId: "my_user",
				},
			},
		},
		{
			Kind: &filepb.Content_Code{
				Code: &filepb.Code{
					Text: "let a = 10",
				},
			},
		},
		{
			Kind: &filepb.Content_Equation{
				Equation: &filepb.Equation{
					Expression: "a + b",
				},
			},
		},
		{
			Kind: &filepb.Content_Image{
				Image: &filepb.Image{
					Url: "https://domain.com/myimage",
				},
			},
		},
		{
			Kind: &filepb.Content_EmbeddedImage{
				EmbeddedImage: &filepb.EmbeddedImage{
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
