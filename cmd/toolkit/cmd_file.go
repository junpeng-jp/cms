package main

import (
	"fmt"
	"os"

	"github.com/junpeng.ong/blog/internal/encoding/storage"
	"github.com/junpeng.ong/blog/internal/file/contentpb"
	"github.com/junpeng.ong/blog/internal/file/metadatapb"
	"github.com/protocolbuffers/protoscope"
	"github.com/spf13/cobra"
)

func newFileComand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "file",
		Short: "utilities to work with a file",
		Long:  ``,
	}

	cmd.AddCommand(
		newFileDecodeCommand(),
	)

	return cmd
}

const fileStringFormat = `---- file start -----
<file start marker>
--- content ---
%s
--- footer ---
%s
(footer size = %d)
<file end marker>
---- file end -----
`

func newFileDecodeCommand() *cobra.Command {
	var fileName string

	cmd := &cobra.Command{
		Use:   "decode",
		Short: "decodes a file into a human readable specification",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := os.ReadFile(fileName)
			if err != nil {
				return err
			}

			bNoMarker, err := storage.StripFileMarkers(b)
			if err != nil {
				return err
			}

			footerSart, footerEnd := storage.ComputeFooterByteRange(bNoMarker)
			footerBytes := bNoMarker[footerSart:footerEnd]
			footer, err := storage.DecodeFooter(footerBytes)
			if err != nil {
				return err
			}

			var contentTotalSize int
			for _, contentMetadata := range footer.ContentRange {
				contentTotalSize += int(contentMetadata.Size)
			}
			contentBytes := bNoMarker[footer.ContentStartOffset : int(footer.ContentStartOffset)+contentTotalSize]

			var contentMsg contentpb.Content
			var footerMsg metadatapb.Footer

			contentRepr := protoscope.Write(
				contentBytes,
				protoscope.WriterOptions{
					NoQuotedStrings:        false,
					AllFieldsAreMessages:   false,
					NoGroups:               true,
					ExplicitWireTypes:      true,
					ExplicitLengthPrefixes: true,
					Schema:                 contentMsg.ProtoReflect().Descriptor(),
					PrintFieldNames:        false,
					PrintEnumNames:         false,
				},
			)

			footerRepr := protoscope.Write(
				footerBytes,
				protoscope.WriterOptions{
					NoQuotedStrings:        false,
					AllFieldsAreMessages:   false,
					NoGroups:               true,
					ExplicitWireTypes:      true,
					ExplicitLengthPrefixes: true,
					Schema:                 footerMsg.ProtoReflect().Descriptor(),
					PrintFieldNames:        true,
					PrintEnumNames:         true,
				},
			)

			fmt.Printf(
				fileStringFormat,
				contentRepr,
				footerRepr,
				footerEnd-footerSart,
			)

			return nil
		},
	}

	cmd.Flags().StringVarP(&fileName, cmdMarkdownFileNameFlagName, cmdMarkdownFileNameFlagShortName, "", "name of content file")

	return cmd
}
