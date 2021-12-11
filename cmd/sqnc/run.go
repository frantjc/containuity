package main

import (
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "run",
}

var (
	jobName     string
	runtimeName string
	stepID      string
)

func init() {
	runCmd.AddCommand(
		runStepCmd,
		runPluginCmd,
	)
}
