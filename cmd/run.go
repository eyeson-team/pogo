package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"pogo/gitlab"
	"pogo/podman"

	"github.com/spf13/cobra"
)

// init adds run command.
func init() {
	rootCmd.AddCommand(runCmd)
}

// runCmd processes the run job of a CI job.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the current CI job",
	Long: `The CI job will run the exec command with an argument that points to a
temporary script file that contains the commands to run. We copy this file
into the running container and run it inside using sh.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		container := podman.NewContainer(gitlab.GetContainerName(), gitlab.GetContainerImage())

		runnerStep := os.Args[len(os.Args)-1]
		scriptFile := os.Args[len(os.Args)-2]
		if Config.Debug {
			CopyScriptForDebug(scriptFile, runnerStep)
		}
		fileName := filepath.Base(scriptFile)
		container.Copy(scriptFile, fileName)
		if err := container.Exec("chmod", "+x", fileName); err != nil {
			return err
		}
		return container.Exec("./" + fileName)
	},
}

// CopyScriptForDebug copies the script for debug to a tempfile. Requires a
// sleep step in the CI job in order to check the file online.
func CopyScriptForDebug(sourceFile, step string) (*os.File, error) {
	tmpfile, err := ioutil.TempFile("", "pogo-ci-debug-"+step)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Copy script file %v to %v\n", string(sourceFile), tmpfile.Name())
	sourceContent, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return nil, err
	}
	_, err = tmpfile.Write(sourceContent)
	return tmpfile, err
}
