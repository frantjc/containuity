package main

import (
	"io"
	"os"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/orchestrator"
	"github.com/frantjc/sequence/runtime"
	"github.com/spf13/cobra"
)

var runWorkflowCmd = &cobra.Command{
	RunE: runRunWorkflow,
	Use:  "workflow",
	Args: cobra.MinimumNArgs(1),
}

func runRunWorkflow(cmd *cobra.Command, args []string) error {
	var (
		ctx  = cmd.Context()
		path = args[0]
		r    io.Reader
		err  error
	)
	if path == fromStdin {
		r = os.Stdin
	} else {
		var err error
		r, err = os.Open(path)
		if err != nil {
			return err
		}
	}

	workflow, err := sequence.NewWorkflowFromReader(r)
	if err != nil {
		return err
	}

	rt, err := runtime.Get(ctx, runtimeName)
	if err != nil {
		return err
	}

	return orchestrator.RunWorkflow(ctx, rt, workflow, orchestrator.WithGitHubToken(gitHubToken))
}
