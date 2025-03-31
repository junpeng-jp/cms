package v0

import (
	"testing"

	"github.com/junpeng.ong/blog/internal/file"
	"github.com/junpeng.ong/blog/internal/file/filepb"
	"github.com/stretchr/testify/assert"
)

func TestStorageCodecIdempotence(t *testing.T) {
	canonicalFile := &file.File{
		ContentList: []*filepb.Content{
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

	testCases := []struct {
		name  string
		codec StorageCodec
		in    *file.File
		want  *file.File
	}{
		{
			name:  "no error: canonical file",
			codec: StorageCodec{},
			in:    canonicalFile,
			want:  canonicalFile,
		},
		{
			name:  "no error: file content list has 1 text",
			codec: StorageCodec{},
			in: &file.File{
				ContentList: []*filepb.Content{
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
				},
			},
			want: &file.File{
				ContentList: []*filepb.Content{
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
				},
			},
		},
		{
			name:  "no error: file content list has 1 text",
			codec: StorageCodec{},
			in: &file.File{
				ContentList: []*filepb.Content{
					{
						Kind: &filepb.Content_Link{
							Link: &filepb.Link{
								Url:  "https://google.com",
								Text: "link to url",
							},
						},
					},
				},
			},
			want: &file.File{
				ContentList: []*filepb.Content{
					{
						Kind: &filepb.Content_Link{
							Link: &filepb.Link{
								Url:  "https://google.com",
								Text: "link to url",
							},
						},
					},
				},
			},
		},
		{
			name:  "no error: file content list is empty array",
			codec: StorageCodec{},
			in: &file.File{
				ContentList: []*filepb.Content{},
			},
			want: &file.File{
				ContentList: []*filepb.Content{},
			},
		},
		{
			name:  "no error: file content list is nil",
			codec: StorageCodec{},
			in: &file.File{
				ContentList: nil,
			},
			want: &file.File{
				ContentList: []*filepb.Content{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.codec.Encode(tc.in)
			assert.NoError(t, err)

			t.Logf("length: %d", len(b))

			decoded, err := tc.codec.Decode(b)
			assert.NoError(t, err)

			assert.Equal(t, tc.want, decoded)
		})
	}
}
