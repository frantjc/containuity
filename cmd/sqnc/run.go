package main

import (
	_ "github.com/frantjc/sequence/runtime/containerd"
	_ "github.com/frantjc/sequence/runtime/docker"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use: "run",
}

var (
	jobName     string
	runtimeName string
	stepID      string
	fromStdin   = "-"
)

func init() {
	runCmd.AddCommand(
		runStepCmd,
	)
}
