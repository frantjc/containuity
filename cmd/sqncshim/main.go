package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/frantjc/sequence/meta"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     fmt.Sprintf("%sshim", meta.Name),
	Version: fmt.Sprintf("%s%s %s", meta.Name, meta.Semver(), runtime.Version()),
}

func init() {
	rootCmd.SetVersionTemplate("{{ .Version }}\n")
	rootCmd.AddCommand(
		pluginCmd,
		versionCmd,
	)
}

func main() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
