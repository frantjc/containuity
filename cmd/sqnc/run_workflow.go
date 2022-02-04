package main

import (
	"context"
	"io"
	"os"

	"github.com/frantjc/sequence"
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

	getConfig()
	workflow, err := sequence.NewWorkflowFromReader(r)
	if err != nil {
		return err
	}

	rt, err := runtime.Get(ctx, runtimeName)
	if err != nil {
		return err
	}

	return runWorkflow(ctx, rt, workflow, withGitHubToken(gitHubToken))
}

func runWorkflow(ctx context.Context, r runtime.Runtime, w *sequence.Workflow, opts ...runOpt) error {
	for name, job := range w.Jobs {
		err := runJob(ctx, r, &job, append(opts, withJobName(name), withJob(&job))...)
		if err != nil {
			return err
		}
	}

	return nil
}
