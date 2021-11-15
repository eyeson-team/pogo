package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// init adds print command.
func init() {
	rootCmd.AddCommand(printCmd)
}

// printCmd prints the configuration.
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print the current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := viper.AllSettings()
		bs, err := yaml.Marshal(c)
		if err != nil {
			return err
		}
		fmt.Println(string(bs))
		return nil
	},
}
