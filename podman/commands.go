package podman

import (
	"os/exec"
)

// Run starts a container.
func Run(containerName, image string) error {
	cmd := exec.Command("podman", "run", "--detach", "--name", containerName, image)
	_, err := cmd.Output()

	return err
}
