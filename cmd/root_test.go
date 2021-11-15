package cmd

import (
	"bytes"
	"os"
	"testing"

	"pogo/config"

	"github.com/spf13/viper"
)

func TestWorkingDirectoryAvailable(t *testing.T) {
	initConfig()

	if Config.WorkingDir != "/home/runner" {
		t.Errorf("Working directory must have default path /home/runner, got %v.", Config.WorkingDir)
	}
}

func TestFileJoinsWithWorkingDir(t *testing.T) {
	initConfig()

	if GetFile("bin/test") != "/home/runner/bin/test" {
		t.Error("File should join a filename with working directory.")
	}
}

func TestConfigurationWorkingDir(t *testing.T) {
	var c config.Config
	yaml := []byte(`working_dir: /home/pogo`)
	viper.ReadConfig(bytes.NewBuffer(yaml))
	err := viper.Unmarshal(&c)
	if err != nil {
		t.Error("Config decoding should not have errors:", err)
	}
	if c.WorkingDir != "/home/pogo" {
		t.Errorf("Could not read configuration working directory, got %v.", c.WorkingDir)
	}
}

func TestConfigurationAuthFile(t *testing.T) {
	var c config.Config
	yaml := []byte(`auth_file: /home/pogo/.auth.json`)
	viper.ReadConfig(bytes.NewBuffer(yaml))
	err := viper.Unmarshal(&c)
	if err != nil {
		t.Error("Config decoding should not have errors:", err)
	}
	if c.AuthFile != "/home/pogo/.auth.json" {
		t.Errorf("Could not read configuration auth file, got %v.", c.AuthFile)
	}
}

func TestConfigurationMounts(t *testing.T) {
	var c config.Config
	yaml, err := os.ReadFile("../fixtures/pogo.yaml")
	if err != nil {
		t.Error("Reading fixtures config files should not have errors:", err)
	}
	viper.ReadConfig(bytes.NewBuffer(yaml))
	err = viper.Unmarshal(&c)
	if err != nil {
		t.Error("Config decoding should not have errors:", err)
	}
	if len(c.Mounts) != 2 {
		t.Errorf("Test config must have 2 mounts, got %v", len(c.Mounts))
	}
}

func TestConfigurationExtraArguments(t *testing.T) {
	var c config.Config
	yaml, err := os.ReadFile("../fixtures/pogo.yaml")
	if err != nil {
		t.Error("Reading fixtures config files should not have errors:", err)
	}
	viper.ReadConfig(bytes.NewBuffer(yaml))
	err = viper.Unmarshal(&c)
	if err != nil {
		t.Error("Config decoding should not have errors:", err)
	}
	if c.Arguments["help"] != "--help" {
		t.Errorf("Extra arguments do miss expected configuration, got %v", c.Arguments)
	}
}
