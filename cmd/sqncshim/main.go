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
	Use:     "sqncshim",
	Version: fmt.Sprintf("%s%s %s", meta.Name, meta.Semver(), runtime.Version()),
}

func init() {
	rootCmd.SetVersionTemplate(
		fmt.Sprintf("{{ with .Version }}{{ . }}{{ end }}\n"),
	)

	rootCmd.AddCommand(
		pluginCmd,
	)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
