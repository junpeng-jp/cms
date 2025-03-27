package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

const (
	cmdMarkdownFileNameFlagName      = "file"
	cmdMarkdownFileNameFlagShortName = "f"
)

func newMarkdownCommand() *cobra.Command {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	cmd := &cobra.Command{
		Use:   "markdown",
		Short: "markdown related utilities",
		Long:  ``,
	}

	cmd.AddCommand(
		newMarkdownParseCommand(md),
	)

	return cmd
}

func newMarkdownParseCommand(md goldmark.Markdown) *cobra.Command {
	var fileName string

	cmd := &cobra.Command{
		Use:   "parse",
		Short: "parses markdown using goldmark",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := os.ReadFile(fileName)
			if err != nil {
				return err
			}
			return parseMardown(md.Parser(), b)
		},
	}

	cmd.Flags().StringVarP(&fileName, cmdMarkdownFileNameFlagName, cmdMarkdownFileNameFlagShortName, "", "name of content file")

	return cmd
}

func parseMardown(parser parser.Parser, markdownBytes []byte) error {
	reader := text.NewReader(markdownBytes)
	doc := parser.Parse(reader)

	doc.Dump(markdownBytes, 1)
	return nil
}
