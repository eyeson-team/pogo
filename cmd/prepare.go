package cmd

import (
	"errors"
	"pogo/gitlab"
	"pogo/podman"

	"github.com/spf13/cobra"
)

// init adds prepare command.
func init() {
	rootCmd.AddCommand(prepareCmd)
}

// prepareCmd runs the prepare step of a CI job.
var prepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Prepare the current CI job for execution",
	RunE: func(cmd *cobra.Command, args []string) error {
		image, err := containerImage()
		if err != nil {
			return err
		}
		container := podman.NewContainer(gitlab.GetContainerName(), image)
		if container.Exists() {
			container.Remove()
		}

		if err = container.Run(&Config); err != nil {
			return err
		}

		if err = startServices(container.Name); err != nil {
			return err
		}

		if err = container.Copy(GetFile("bin/gitlab-runner"), "/usr/bin/gitlab-runner"); err != nil {
			return err
		}
		if err = container.Copy(GetFile("bin/install-packages"), "/usr/bin/install-packages"); err != nil {
			return err
		}

		if err = container.Exec("install-packages", "git", "ca-certificates"); err != nil {
			return err
		}

		return nil
	},
}

// startServices starts all required job services and runs them within the
// network of the job container.
func startServices(containerName string) error {
	serviceImages, err := gitlab.GetServiceImages()
	if err != nil {
		return err
	}
	for _, serviceImage := range serviceImages {
		serviceContainer := podman.NewContainer(gitlab.GetServiceContainerName(serviceImage.Name),
			serviceImage.Name)
		if err = serviceContainer.RunFor(containerName, &Config); err != nil {
			return err
		}
	}
	return nil
}

// containerImage provides the container image from CI environment variables or
// falls back to the default image defined in config.
func containerImage() (string, error) {
	image := gitlab.GetContainerImage()
	if image == "" {
		image = Config.DefaultImage
	}
	if image == "" {
		return "", errors.New("No job- or default container image set.")
	}
	return image, nil
}
