package main

import (
	"io"
	"os"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/orchestrator"
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

	return orchestrator.RunJob(ctx, rt, job, orchestrator.WithGitHubToken(gitHubToken), orchestrator.WithJobName(jobName))
}
