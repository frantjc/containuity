package main

import (
	"github.com/frantjc/sequence/defaults"
	_ "github.com/frantjc/sequence/runtime/containerd"
	_ "github.com/frantjc/sequence/runtime/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:               "run",
	PersistentPreRunE: runPersistentPreRun,
}

const (
	fromStdin = "-"
)

var (
	runtimeName string
)

var (
	jobName string
	stepID  string
)

func init() {
	runStepCmd.PersistentFlags().StringVarP(&runtimeName, "runtime", "", defaults.Runtime, "container runtime to use")
	runCmd.AddCommand(
		runStepCmd,
	)
	viper.BindPFlags(runCmd.Flags())
}

func runPersistentPreRun(cmd *cobra.Command, args []string) error {
	viper.SafeWriteConfig()
	return nil
}
