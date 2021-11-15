package cmd

import (
	"fmt"
	"pogo/gitlab"
	"pogo/podman"

	"github.com/spf13/cobra"
)

// init adds clean command.
func init() {
	rootCmd.AddCommand(cleanCmd)
}

// cleanCmd removes running service- and job containers.
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the current CI job after execution",
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceImages, err := gitlab.GetServiceImages()
		if err != nil {
			return err
		}
		fmt.Println(len(serviceImages))
		for _, serviceImage := range serviceImages {
			serviceContainer := podman.NewContainer(gitlab.GetServiceContainerName(serviceImage.Name),
				serviceImage.Name)
			err = serviceContainer.Remove()
			if err != nil {
				return err
			}
		}
		container := podman.NewContainer(gitlab.GetContainerName(), "")
		return container.Remove()
	},
}
