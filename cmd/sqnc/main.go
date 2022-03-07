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
	rootCmd.SetVersionTemplate(
		fmt.Sprintf("{{ with .Name }}{{ . }}{{ end }}{{ with .Version }}{{ . }}{{ end }} %s\n", runtime.Version()),
	)

	rootCmd.PersistentFlags().Bool("verbose", false, "verbose")
	rootCmd.PersistentFlags().Int("port", 0, "port")
	rootCmd.PersistentFlags().String("sock", "", "unix socket")

	conf.BindVerboseFlag(rootCmd.Flag("verbose"))
	conf.BindPortFlag(rootCmd.Flag("port"))
	conf.BindSocketFlag(rootCmd.Flag("socket"))
	// conf.BindGitHubTokenFlag(rootCmd.Flag("github-token"))
	// conf.BindRuntimeImageFlag(rootCmd.Flag("runtime-image"))
	// conf.BindRuntimeNameFlag(rootCmd.Flag("runtime-name"))
	// conf.BindSecretsFlag(rootCmd.Flag("secret"))

	rootCmd.AddCommand(
		runCmd,
	)
}

func persistentPreRun(cmd *cobra.Command, args []string) error {
	c, err := conf.Get()
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
