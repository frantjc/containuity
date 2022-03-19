package main

import "github.com/spf13/cobra"

const fromStdin = "-"

var (
	runCmd = &cobra.Command{
		Use:  "run",
		RunE: runWorkflow,
	}
	jobName string
	stepID  string
)

func init() {
	runCmd.AddCommand(
		runStepCmd,
		runJobCmd,
		runWorkflowCmd,
	)
}
