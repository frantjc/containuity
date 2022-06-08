package main

import (
	"io"
	"os"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/conf"
	"github.com/frantjc/sequence/internal/conf/flags"
	"github.com/frantjc/sequence/internal/log"
	workflowv1 "github.com/frantjc/sequence/workflow/v1"
	"github.com/spf13/cobra"
)

var runJobCmd = &cobra.Command{
	Use:  "job",
	Args: cobra.ExactArgs(1),
	RunE: runJob,
}

func init() {
	runJobCmd.Flags().StringVar(&jobName, "job", "", "name of the job to run")
}

func runJob(cmd *cobra.Command, args []string) error {
	var (
		ctx  = cmd.Context()
		j    *workflowv1.Job
		w    *workflowv1.Workflow
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

	if jobName != "" {
		w, err := workflowv1.NewWorkflowFromReader(r)
		if err != nil {
			return err
		}

		j, err = w.GetJob(jobName)
		if err != nil {
			return err
		}
	} else {
		j, err = workflowv1.NewJobFromReader(r)
		if err != nil {
			return err
		}
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
		sequence.WithWorkflow(w),
		sequence.WithRunnerImage(c.Runtime.RunnerImage),
		sequence.WithRepository(flags.FlagWorkDir),
	}
	if c.Verbose {
		opts = append(opts, sequence.WithVerbose)
	}

	return client.RunJob(ctx, j, log.Writer(), opts...)
}
