package cmd

import (
	"fmt"
	"os"
	"pogo/config"
	"pogo/gitlab"

	"encoding/json"

	"github.com/spf13/cobra"
)

type configOutputDriver struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type configOutput struct {
	BuildsDir string             `json:"builds_dir"`
	CacheDir  string             `json:"cache_dir"`
	Hostname  string             `json:"hostname"`
	Driver    configOutputDriver `json:"driver"`
	JobEnv    map[string]string  `json:"job_env"`
}

// init adds config command.
func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the current CI job",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := &configOutput{
			BuildsDir: gitlab.GetBuildsDir(),
			CacheDir:  gitlab.GetCacheDir(),
			Hostname:  os.Getenv("HOSTNAME"),
			Driver: configOutputDriver{
				Name:    "pogo",
				Version: config.VERSION,
			},
		}
		jsonOut, err := json.Marshal(out)
		if err != nil {
			return err
		}
		fmt.Println(string(jsonOut))
		return nil
	},
}
