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
	"github.com/frantjc/sequence/conf"
	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/meta"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               fmt.Sprintf("%sd", meta.Name),
	Version:           meta.Semver(),
	PersistentPreRunE: persistentPreRun,
	RunE:              run,
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
}

func persistentPreRun(cmd *cobra.Command, args []string) error {
	c, err := conf.NewFull()
	if err != nil {
		return err
	}
	log.SetVerbose(c.Verbose)
	return nil
}

func run(cmd *cobra.Command, args []string) error {
	var (
		ctx    = cmd.Context()
		c, err = conf.NewFull()
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

	s, err := sequence.NewServer(ctx, sequence.WithAnyRuntime)
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
