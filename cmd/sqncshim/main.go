package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/frantjc/sequence"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "sqncshim",
	Version: fmt.Sprintf("sqncshim%s %s", sequence.Semver(), runtime.Version()),
}

func init() {
	rootCmd.SetVersionTemplate("{{ .Version }}\n")
	rootCmd.AddCommand(
		pluginCmd,
		stepCmd,
		versionCmd,
	)
}

func main() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
