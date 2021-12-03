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
	SilenceErrors:    true,
	SilenceUsage:     true,
	TraverseChildren: true,
	Use:              sequence.Name,
	Version:          sequence.Version,
}

func init() {
	rootCmd.AddCommand(
		runCmd,
	)
	rootCmd.SetVersionTemplate(
		fmt.Sprintf("{{ with .Name }}{{ . }}{{ end }}{{ .Version }} %s\n", runtime.Version()),
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
