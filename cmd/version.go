package cmd

import (
	"fmt"
	"pogo/config"

	"github.com/spf13/cobra"
)

// init adds config command.
func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version of pogo",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("pogo %v\n", config.VERSION)
		return nil
	},
}
