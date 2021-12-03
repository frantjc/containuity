package main

import (
	"context"
	"encoding/base64"

	"github.com/frantjc/sequence/key"
	_ "github.com/frantjc/sequence/runtime"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "run",
}

var (
	jobName     string
	runtimeName string
	stepID      string
)

func init() {
	runCmd.AddCommand(
		runStepCmd,
	)
}

func encode(ss ...string) string {
	var sa string
	for _, s := range ss {
		sa += s
	}
	return base64.RawStdEncoding.EncodeToString([]byte(sa))
}

func withJob(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, key.Job, id)
}
