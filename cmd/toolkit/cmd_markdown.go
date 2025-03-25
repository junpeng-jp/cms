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

func newGenCommand() *cobra.Command {
	var fileName string

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

	cmdMarkdown := &cobra.Command{
		Use:   "markdown",
		Short: "markdown related utilities",
		Long:  ``,
	}

	cmdMarkdownParse := &cobra.Command{
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

	cmdMarkdown.AddCommand(
		cmdMarkdownParse,
	)

	cmdMarkdownParse.Flags().StringVarP(&fileName, cmdMarkdownFileNameFlagName, cmdMarkdownFileNameFlagShortName, "", "name of content file")

	return cmdMarkdown
}

func parseMardown(parser parser.Parser, markdownBytes []byte) error {
	reader := text.NewReader(markdownBytes)
	doc := parser.Parse(reader)

	doc.Dump(markdownBytes, 1)
	return nil
}
