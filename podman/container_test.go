package podman

import (
	"pogo/config"
	"strings"
	"testing"
)

func TestGetServiceImagesCanHandleEmptyValue(t *testing.T) {
	container := NewContainer("test-container", "test-image")

	if container.Exists() == true {
		t.Error("Test container is not expected to exist.")
	}
}

func TestAddMount(t *testing.T) {
	mount := config.VolumeMountConfig{"/home:/mnt/home", []string{}, true}
	opts := []string{"run"}
	addMount(&opts, &mount)
	if strings.Join(opts, ",") != "run,-v,/home:/mnt/home:z,ro" {
		t.Error("Failed to receive expected volume opts, got:", opts)
	}
}
