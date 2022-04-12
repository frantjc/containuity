package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/frantjc/sequence"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:  "version",
	RunE: runVersion,
}

func runVersion(cmd *cobra.Command, args []string) error {
	return write(os.Stdout, fmt.Sprintf("sqnc%s %s", sequence.Semver(), runtime.Version()))
}
