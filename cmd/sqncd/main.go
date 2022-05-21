package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/conf"
	"github.com/frantjc/sequence/internal/conf/flags"
	"github.com/frantjc/sequence/internal/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "sqncd",
	Version:           sequence.Semver(),
	PersistentPreRunE: persistentPreRun,
	RunE:              run,
}

func init() {
	rootCmd.SetVersionTemplate(
		fmt.Sprintf("{{ .Name }}{{ .Version }} %s\n", runtime.Version()),
	)
	rootCmd.PersistentFlags().StringVar(&flags.FlagConfigFilePath, "config", "", "config file")
	rootCmd.PersistentFlags().BoolVar(&flags.FlagVerbose, "verbose", false, "verbose")
	rootCmd.PersistentFlags().StringVar(&flags.FlagSocket, "sock", "", "unix socket")
	rootCmd.PersistentFlags().Int64Var(&flags.FlagPort, "port", 0, "port")
	rootCmd.PersistentFlags().StringVar(&flags.FlagRootDir, "root-dir", "", "root dir")
	rootCmd.PersistentFlags().StringVar(&flags.FlagStateDir, "state-dir", "", "state dir")
	wd, _ := os.Getwd()
	rootCmd.PersistentFlags().StringVar(&flags.FlagWorkDir, "context", wd, "context")
	rootCmd.AddCommand(versionCmd)
}

func persistentPreRun(cmd *cobra.Command, args []string) error {
	c, err := conf.NewFromFlags()
	if err != nil {
		return err
	}
	log.SetVerbose(c.Verbose)
	return nil
}

func run(cmd *cobra.Command, args []string) error {
	var (
		ctx    = cmd.Context()
		c, err = conf.NewFromFlags()
	)
	if err != nil {
		return err
	}

	addr := strings.TrimPrefix(c.Address(), "unix://")
	os.MkdirAll(c.RootDir, 0777)
	os.MkdirAll(c.StateDir, 0777)
	if c.Port == 0 {
		os.MkdirAll(filepath.Dir(addr), 0777)
	}
	os.Remove(addr)
	defer os.Remove(addr)

	l, err := net.Listen(c.Network(), addr)
	if err != nil {
		return err
	}

	s, err := sequence.NewServer(ctx, sequence.WithRuntimeName(c.Runtime.Name))
	if err != nil {
		return err
	}

	return s.Serve(l)
}

func main() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
