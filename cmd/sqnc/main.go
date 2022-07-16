package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/flags"
	"github.com/frantjc/sequence/internal/plugins"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "sqnc",
		Version: sequence.Semver(),
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			if err := plugins.Open(); err != nil {
				cmd.PrintErrln(err)
			}
		},
	}
)

func init() {
	rootCmd.SetVersionTemplate(
		fmt.Sprintf("{{ .Name }}{{ .Version }} %s\n", runtime.Version()),
	)
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flags.PluginDir, "plugins", "", "plugin directory")
	// TODO if err := rootCmd.MarkFlagDirname("plugins"); err != nil {
	// 	panic(err)
	// }
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
