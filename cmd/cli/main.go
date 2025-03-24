package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "toolkit",
	Short: "toolkit is a content management CLI",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(
		GenCommand,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
