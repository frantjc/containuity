package main

import (
	"github.com/frantjc/sequence/defaults"
	_ "github.com/frantjc/sequence/runtime/containerd"
	_ "github.com/frantjc/sequence/runtime/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use: "run",
	PersistentPreRunE: persistentPreRunRun,
}

const (
	fromStdin = "-"
)

var (
	runtimeName string
	jobName     string
	stepID      string
	gitHubToken string
)

func init() {
	runCmd.PersistentFlags().StringP("runtime", "r", defaults.Runtime, "container runtime to use")
	runCmd.PersistentFlags().String("github-token", "", "GitHub token")

	viper.BindPFlag("runtime.name", runCmd.Flag("runtime"))
	viper.BindPFlag("github.token", runCmd.Flag("github-token"))

	runCmd.AddCommand(
		runStepCmd,
		runJobCmd,
		runWorkflowCmd,
	)
}

func persistentPreRunRun(cmd *cobra.Command, args []string) error {
	viper.ReadInConfig()
	return nil
}

func getConfig() {
	gitHubToken = viper.GetString("github.token")
	runtimeName = viper.GetString("runtime.name")
}
