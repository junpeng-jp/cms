package storage

import (
	"testing"

	"github.com/junpeng.ong/blog/internal/file"
	"github.com/junpeng.ong/blog/internal/file/contentpb"
	"github.com/stretchr/testify/assert"
)

func TestStorageCodecIdempotence(t *testing.T) {
	canonicalFile := &file.File{
		ContentList: []*contentpb.Content{
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
				ContentList: []*contentpb.Content{
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
				},
			},
			want: &file.File{
				ContentList: []*contentpb.Content{
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
				},
			},
		},
		{
			name:  "no error: file content list has 1 text",
			codec: StorageCodec{},
			in: &file.File{
				ContentList: []*contentpb.Content{
					{
						Kind: &contentpb.Content_Link{
							Link: &contentpb.Link{
								Url:  "https://google.com",
								Text: "link to url",
							},
						},
					},
				},
			},
			want: &file.File{
				ContentList: []*contentpb.Content{
					{
						Kind: &contentpb.Content_Link{
							Link: &contentpb.Link{
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
				ContentList: []*contentpb.Content{},
			},
			want: &file.File{
				ContentList: []*contentpb.Content{},
			},
		},
		{
			name:  "no error: file content list is nil",
			codec: StorageCodec{},
			in: &file.File{
				ContentList: nil,
			},
			want: &file.File{
				ContentList: []*contentpb.Content{},
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
