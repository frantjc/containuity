package main

import (
	"fmt"
	"runtime"

	"github.com/frantjc/sequence/meta"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:  "version",
	RunE: runVersion,
}

func runVersion(cmd *cobra.Command, args []string) error {
	return write(cmd.OutOrStdout(), fmt.Sprintf("%s%s %s", meta.Name, meta.Semver(), runtime.Version()))
}