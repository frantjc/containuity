package main

import (
	"context"
	"io"
	"os"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/runtime"
	"github.com/spf13/cobra"
)

var runJobCmd = &cobra.Command{
	RunE: runRunJob,
	Use:  "job",
	Args: cobra.MinimumNArgs(1),
}

func init() {
	runJobCmd.Flags().StringVarP(&jobName, "job", "j", "", "name of the job to run")
}

func runRunJob(cmd *cobra.Command, args []string) error {
	var (
		ctx  = cmd.Context()
		job  *sequence.Job
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
	if jobName != "" {
		workflow, err := sequence.NewWorkflowFromReader(r)
		if err != nil {
			return err
		}

		job, err = workflow.GetJob(jobName)
		if err != nil {
			return err
		}
	} else {
		job, err = sequence.NewJobFromReader(r)
		if err != nil {
			return err
		}
	}

	rt, err := runtime.Get(ctx, runtimeName)
	if err != nil {
		return err
	}

	return runJob(ctx, rt, job, withGitHubToken(gitHubToken), withJobName(jobName))
}

func runJob(ctx context.Context, r runtime.Runtime, j *sequence.Job, opts ...runOpt) error {
	for _, step := range j.Steps {
		err := runStep(ctx, r, &step, append(opts, withJob(j))...)
		if err != nil {
			return err
		}
	}

	return nil
}
