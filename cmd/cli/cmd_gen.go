package main

import (
	"github.com/spf13/cobra"
)

var GenCommand = &cobra.Command{
	Use:   "gen",
	Short: "generates the relevant static asets",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {

}
