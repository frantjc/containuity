package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/meta"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     meta.Name,
	Version: meta.Semver(),
	PersistentPreRunE: persistentPreRun,
}

const (
	configName = "config"
)

var (
	home    string
	verbose bool
)

func init() {
	home = os.Getenv("HOME")
	if home == "" {
		home = "$HOME"
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
	rootCmd.AddCommand(
		runCmd,
		pluginCmd,
	)
	rootCmd.SetVersionTemplate(
		fmt.Sprintf("{{ with .Name }}{{ . }}{{ end }}{{ with .Version }}{{ . }}{{ end }} %s\n", runtime.Version()),
	)
	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(fmt.Sprintf("%s/.%s", home, meta.Name))
	viper.AddConfigPath(fmt.Sprintf("/etc/%s", meta.Name))
	viper.SetEnvPrefix(meta.Name)
	viper.AllowEmptyEnv(true)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

func persistentPreRun(cmd *cobra.Command, args []string) error {
	log.SetVerbose(verbose)
	return viper.ReadInConfig()
}
