package main

import (
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use: "plugin",
}

func init() {
	pluginCmd.AddCommand(
		pluginUsesCmd,
	)
}
