package main

import (
	"github.com/frantjc/sequence/conf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:               "run",
	PersistentPreRunE: persistentPreRunRun,
	RunE:              runWorkflow,
}

const (
	fromStdin = "-"
)

var (
	jobName string
	stepID  string
)

func init() {
	runCmd.PersistentFlags().StringP("runtime", "r", "", "container runtime to use")
	runCmd.PersistentFlags().StringP("image", "i", "", "default image to use")
	runCmd.PersistentFlags().String("github-token", "", "GitHub token")

	conf.BindRuntimeNameFlag(runCmd.Flag("runtime"))
	conf.BindGitHubTokenFlag(runCmd.Flag("image"))
	conf.BindGitHubTokenFlag(runCmd.Flag("github-token"))

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
