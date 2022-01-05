package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"

	"github.com/frantjc/sequence"
	_ "github.com/frantjc/sequence/pkg/runtime/docker"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var rootCmd = &cobra.Command{
	Use:               sequence.Name,
	Version:           sequence.Version,
	PersistentPreRunE: persistentPreRun,
	RunE:              run,
}

var (
	verbose bool
	socket  string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	rootCmd.PersistentFlags().StringVarP(&socket, "socket", "s", "/tmp/sequence.sock", "Unix socket")
	rootCmd.AddCommand(
		runCmd,
		pluginCmd,
	)
	rootCmd.SetVersionTemplate(
		fmt.Sprintf("{{ with .Name }}{{ . }}{{ end }}{{ with .Version }}{{ . }}{{ end }} %s\n", runtime.Version()),
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

func persistentPreRun(cmd *cobra.Command, args []string) error {
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return nil
}

func run(cmd *cobra.Command, args []string) error {
	socket = strings.TrimPrefix(socket, "unix://")
	l, err := net.Listen("unix", socket)
	if err != nil {
		return err
	}
	defer l.Close()

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	defer s.GracefulStop()

	return s.Serve(l)
}
