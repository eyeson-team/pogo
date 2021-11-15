package cmd

import (
	"fmt"
	"os"
	"strconv"

	"path/filepath"
	"pogo/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cfgFile provides path of the configuration file.
var cfgFile string

// rootCmd registers the root command that is the base for all other commands.
// It ensures the proper exit status for the job is returned on failure and
// the error message is printed properly.
var rootCmd = &cobra.Command{
	Use:   "pogo",
	Short: "Pogo is a podman gitlab executor",
	Long: `Handling a GitLab custom executor is cumbersome. This friendly helper
does take care for preparing a CI job by authenticating the registry, pulling
the image, fetching the source, starting any service container and cleaning
up stuff afterwards.`,
}

// Config holds the configuration.
var Config config.Config

// init adds root command and registers base flags.
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.pogo.yaml)")
	rootCmd.PersistentFlags().BoolVar(&Config.Debug, "debug", Config.Debug, "print debug output")
	rootCmd.PersistentFlags().StringVarP(&Config.WorkingDir, "working-dir", "w", Config.WorkingDir, "working directory (default: /home/runner)")
}

// GetFile joins a filename with the working directory.
func GetFile(filename string) string {
	return filepath.Join(Config.WorkingDir, filename)
}

// Exit handles to use proper exit code that is provided by GitLab.
func Exit() {
	exitCode, exists := os.LookupEnv("BUILD_FAILURE_EXIT_CODE")
	exitVal, err := strconv.Atoi(exitCode)
	if err != nil || exists == false {
		exitVal = 1
	}
	os.Exit(exitVal)
}

// Execute processes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		Exit()
	}
}

// initConfig initializes configuration.
func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetDefault("debug", false)
	viper.SetDefault("working_dir", "/home/runner")
	viper.SetDefault("cache_dir", "/home/runner/cache")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".pogo")
	}
	viper.AutomaticEnv()
	viper.ReadInConfig()
	if err := viper.Unmarshal(&Config); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode configuration, %v", err)
	}
}
