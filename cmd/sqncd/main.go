package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/conf"
	"github.com/frantjc/sequence/internal/conf/flags"
	// I have no idea what 'File is not `goimports`-ed' means
	"github.com/frantjc/sequence/internal/log"
	// Init default runtime
	_ "github.com/frantjc/sequence/runtime/docker"
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

	if err = os.MkdirAll(c.RootDir, 0777); err != nil {
		return err
	}

	if err = os.MkdirAll(c.StateDir, 0777); err != nil {
		return err
	}

	if c.Port == 0 {
		if err = os.MkdirAll(filepath.Dir(addr), 0777); err != nil {
			return err
		}
	}

	os.Remove(addr)
	defer os.Remove(addr)

	l, err := net.Listen(c.Network(), addr)
	if err != nil {
		return err
	}
	defer l.Close()

	hl, err := net.Listen("tcp", c.HTTPAddress())
	if err != nil {
		return err
	}
	defer hl.Close()

	runtime, err := sequence.GetRuntime(c.Runtime.Name)
	if err != nil {
		return err
	}

	s, err := sequence.NewServer(ctx, runtime)
	if err != nil {
		return err
	}

	errC := make(chan error, 1)

	go func() {
		errC <- s.ServeGRPC(l)
	}()
	log.Infof("gRPC listening on '%s'", l.Addr().String())

	go func() {
		errC <- http.Serve(hl, s)
	}()
	log.Infof("http listening on '%s'", hl.Addr().String())

	return <-errC
}

func main() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
