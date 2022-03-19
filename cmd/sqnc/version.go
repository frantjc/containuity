package main

import "github.com/spf13/cobra"

var versionCmd = &cobra.Command{
	Use:  "version",
	RunE: runVersion,
}

func runVersion(cmd *cobra.Command, args []string) error {
	return write(cmd.OutOrStdout(), rootCmd.Version)
}
