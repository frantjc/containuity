package main

import (
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use: "run",
}

var (
	jobName     string
	runtimeName string
	stepID      string
)

func init() {
	runCmd.AddCommand(
		runStepCmd,
	)
}
