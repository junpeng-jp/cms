package main

import (
	"github.com/spf13/cobra"
)

func newFileComand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "file",
		Short: "utilities to render a file",
		Long:  ``,
	}

	cmd.AddCommand(
		newFileRenderCommand(),
	)

	return cmd
}

func newFileRenderCommand() *cobra.Command {
	var fileName string

	cmd := &cobra.Command{
		Use:   "decode",
		Short: "decodes a file into a human readable specification",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().StringVarP(&fileName, cmdMarkdownFileNameFlagName, cmdMarkdownFileNameFlagShortName, "", "name of content file")

	return cmd
}
