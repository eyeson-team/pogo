package podman

import (
	"bufio"
	"fmt"
	"os/exec"
	"path"
	"pogo/config"
	"pogo/gitlab"
	"strings"
)

// SLEEP_TIMEOUT defines a hard limit a job can run.
const SLEEP_TIMEOUT = "7200"

// container defines a data structure for a single container.
type container struct {
	Name  string
	Image string
}

// Exists checks for a container by inspecting its configuration. In case an
// error is returned, there container does not exist.
func (c *container) Exists() bool {
	cmd := exec.Command("podman", "inspect", c.Name)
	_, err := cmd.Output()

	return err == nil
}

// Remove the container.
func (c *container) Remove() error {
	if err := c.exec("stop", c.Name); err != nil {
		return err
	}

	return c.exec("rm", c.Name)
}

// pull fetches a container image.
func (c *container) pull(authFile string) error {
	if authFile != "" {
		return c.exec("pull", "--authfile", authFile, c.Image)
	}
	return c.exec("pull", c.Image)
}

// Copy copies a file from local into the running container.
func (c *container) Copy(localFile, remoteFile string) error {
	return c.exec("cp", "--pause=false", localFile, c.Name+":"+remoteFile)
}

// addMount adds mount options to podman arguments.
func addMount(opts *[]string, mount *config.VolumeMountConfig) {
	mountOpts := mount.Volume + ":z"
	if mount.Readonly {
		mountOpts = mountOpts + ",ro"
	}
	*opts = append(*opts, "-v", mountOpts)
}

// addCacheDir adds cache directory volume mount to podman arguments.
func addCacheDir(opts *[]string, cacheDir *string) error {
	localDir := path.Join(*cacheDir, gitlab.GetPathSlug())
	if _, err := exec.Command("mkdir", "-p", localDir).Output(); err != nil {
		return err
	}
	remoteDir := gitlab.GetCacheDir()
	*opts = append(*opts, "-v", localDir+":"+remoteDir+":z,rw")
	return nil
}

// Run starts a container and sets it in a passive (sleep) mode.
func (c *container) Run(cfg *config.Config) error {
	if err := c.pull(cfg.AuthFile); err != nil {
		return err
	}
	opts := []string{"-i", "--name", c.Name, "--detach"}
	for tag, args := range cfg.Arguments {
		res, err := gitlab.JobTagsMatch(&[]string{tag})
		if err != nil {
			return err
		}
		if res == true {
			opts = append(opts, strings.Split(args, " ")...)
		}
	}
	for _, mount := range cfg.Mounts {
		if len(mount.Tags) == 0 {
			addMount(&opts, &mount)
		}
		res, err := gitlab.JobTagsMatch(&mount.Tags)
		if err != nil {
			return err
		}
		if res == true {
			addMount(&opts, &mount)
		}
	}
	if err := addCacheDir(&opts, &cfg.CacheDir); err != nil {
		return err
	}
	opts = append([]string{"run"}, opts...)
	opts = append(opts, []string{c.Image, "sleep", SLEEP_TIMEOUT}...)
	return c.exec(opts...)
}

// RunFor starts a container, the given host container given by the name
// arrgument is the target network the new one is attached to.
func (c *container) RunFor(jobContainerName string, cfg *config.Config) error {
	if err := c.pull(cfg.AuthFile); err != nil {
		return err
	}
	return c.exec("run", "-i", "--net", "container:"+jobContainerName, "--name",
		c.Name, "--detach", c.Image)
}

// Exec does execute a command at the given container.
func (c *container) Exec(cmd ...string) error {
	return c.exec(append([]string{"exec", "-i", c.Name}, cmd...)...)
}

// exec executes podman commands prints all as streaoutput.
func (c *container) exec(args ...string) error {
	fmt.Println("$ podman", strings.Join(args, " "))
	cmd := exec.Command("podman", args...)
	reader, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	cmd.Start()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	return cmd.Wait()
}

// NewContainer returns a resource type container.
func NewContainer(name, image string) container {
	return container{name, image}
}
