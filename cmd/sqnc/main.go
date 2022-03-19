package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/frantjc/sequence/conf"
	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/meta"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "sqnc",
	Version:           meta.Semver(),
	PersistentPreRunE: persistentPreRun,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&conf.ConfigFilePath, "config", "", "config file")
	rootCmd.PersistentFlags().BoolVar(&conf.Verbose, "verbose", false, "verbose")
	rootCmd.PersistentFlags().StringVar(&conf.Socket, "sock", "", "unix socket")
	rootCmd.PersistentFlags().IntVar(&conf.Port, "port", 0, "port")
	rootCmd.PersistentFlags().StringVar(&conf.RootDir, "root-dir", "", "root dir")
	rootCmd.PersistentFlags().StringVar(&conf.StateDir, "state-dir", "", "state dir")
	wd, _ := os.Getwd()
	rootCmd.PersistentFlags().StringVar(&conf.WorkDir, "context", wd, "context")
}

func init() {
	rootCmd.SetVersionTemplate(
		fmt.Sprintf("{{ with .Name }}{{ . }}{{ end }}{{ with .Version }}{{ . }}{{ end }} %s\n", runtime.Version()),
	)

	rootCmd.AddCommand(
		runCmd,
		configCmd,
		versionCmd,
	)
}

func persistentPreRun(cmd *cobra.Command, args []string) error {
	c, err := conf.NewFull()
	if err != nil {
		return err
	}
	log.SetVerbose(c.Verbose)
	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
