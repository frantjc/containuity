package main

import (
	"github.com/spf13/cobra"
)

var runPluginCmd = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "plugin",
}

func init() {
	runCmd.AddCommand(
		runPluginUsesCmd,
	)
}
