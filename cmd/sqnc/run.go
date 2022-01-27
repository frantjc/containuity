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
	runCmd.PersistentFlags().StringVarP(&runtimeName, "runtime", "r", defaults.Runtime, "container runtime to use")
	runCmd.PersistentFlags().StringVar(&gitHubToken, "github-token", "", "GitHub token")

	viper.BindPFlag("runtime.name", runCmd.Flag("runtime"))
	viper.BindPFlag("github.token", runCmd.Flag("github-token"))

	runCmd.AddCommand(
		runStepCmd,
		runJobCmd,
		runWorkflowCmd,
	)
}
