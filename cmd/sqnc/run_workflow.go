package main

import (
	"io"
	"os"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/conf"
	"github.com/frantjc/sequence/conf/flags"
	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/workflow"
	"github.com/spf13/cobra"
)

var runWorkflowCmd = &cobra.Command{
	Use:  "workflow",
	Args: cobra.ExactArgs(1),
	RunE: runWorkflow,
}

func runWorkflow(cmd *cobra.Command, args []string) error {
	var (
		ctx  = cmd.Context()
		path = args[0]
		r    io.Reader
		err  error
	)
	if path == fromStdin {
		r = os.Stdin
	} else {
		r, err = os.Open(path)
		if err != nil {
			return err
		}
	}

	w, err := workflow.NewWorkflowFromReader(r)
	if err != nil {
		return err
	}

	c, err := conf.NewFromFlags()
	if err != nil {
		return err
	}

	client, err := sequence.New(ctx, c.Address())
	if err != nil {
		return err
	}

	opts := []sequence.RunOpt{
		sequence.WithRunnerImage(c.Runtime.RunnerImage),
		sequence.WithRepository(flags.FlagWorkDir),
	}
	if c.Verbose {
		opts = append(opts, sequence.WithVerbose)
	}

	return client.RunWorkflow(ctx, w, log.Writer(), opts...)
}
