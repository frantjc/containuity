package main

import (
	"github.com/spf13/cobra"
)

var runPluginCmd = &cobra.Command{
	Use: "plugin",
}

func init() {
	runPluginCmd.AddCommand(
		runPluginUsesCmd,
	)
}
