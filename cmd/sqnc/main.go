package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/frantjc/sequence"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "sqnc",
		Version: sequence.Semver(),
	}
)

func init() {
	rootCmd.SetVersionTemplate(
		fmt.Sprintf("{{ .Name }}{{ .Version }} %s\n", runtime.Version()),
	)
}

func init() {
	rootCmd.AddCommand(
		runCmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
