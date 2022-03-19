package main

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf"},
	RunE:    runConfigGet,
}

func init() {
	configCmd.AddCommand(
		configExampleCmd,
		configGetCmd,
	)
}
