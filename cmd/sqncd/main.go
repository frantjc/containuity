package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"

	"github.com/frantjc/sequence/conf"
	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/meta"
	"github.com/frantjc/sequence/service"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               fmt.Sprintf("%sd", meta.Name),
	Version:           meta.Semver(),
	PersistentPreRunE: persistentPreRun,
	RunE:              run,
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
	conf.BindSocketFlag(rootCmd.Flag("sock"))
	// conf.BindGitHubTokenFlag(rootCmd.Flag("github-token"))
	// conf.BindRuntimeImageFlag(rootCmd.Flag("runtime-image"))
	// conf.BindRuntimeNameFlag(rootCmd.Flag("runtime-name"))
	// conf.BindSecretsFlag(rootCmd.Flag("secret"))
}

func persistentPreRun(cmd *cobra.Command, args []string) error {
	c, err := conf.Get()
	if err != nil {
		return err
	}
	log.SetVerbose(c.Verbose)
	return nil
}

func run(cmd *cobra.Command, args []string) error {
	var (
		ctx    = cmd.Context()
		c, err = conf.Get()
	)
	if err != nil {
		return err
	}

	addr := strings.TrimPrefix(c.Address, "unix://")

	os.Remove(addr)
	defer os.Remove(addr)

	l, err := net.Listen(c.Network, addr)
	if err != nil {
		return err
	}

	s, err := service.New(ctx, service.WithConfig(c))
	if err != nil {
		return err
	}

	return s.Serve(l)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
