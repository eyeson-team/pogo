package cmd

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestDefaultImage(t *testing.T) {
	viper.SetConfigFile("../fixtures/pogo.yaml")
	viper.ReadInConfig()
	if err := viper.Unmarshal(&Config); err != nil {
		t.Errorf("Unable to decode configuration, %v", err)
	}

	image, err := containerImage()
	if err != nil {
		t.Errorf("Failed to read container image, %v", err)
	}
	if image != "fedora:latest" {
		t.Error("Default image was not read from configuration")
	}
	viper.Reset()
}

func TestContainerImage(t *testing.T) {
	viper.SetConfigFile("../fixtures/pogo.yaml")
	viper.ReadInConfig()
	if err := viper.Unmarshal(&Config); err != nil {
		t.Errorf("Unable to decode configuration, %v", err)
	}

	os.Setenv("CUSTOM_ENV_CI_JOB_IMAGE", "alpine:latest")
	image, err := containerImage()
	if err != nil {
		t.Errorf("Failed to read container image, %v", err)
	}
	if image != "alpine:latest" {
		t.Error("Default image was not read from configuration")
	}
	viper.Reset()
}
