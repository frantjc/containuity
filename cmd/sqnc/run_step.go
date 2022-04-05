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

var runStepCmd = &cobra.Command{
	Use:  "step",
	Args: cobra.ExactArgs(1),
	RunE: runStep,
}

func init() {
	runStepCmd.Flags().StringVar(&stepID, "step", "", "ID of the step to run")
	runStepCmd.Flags().StringVar(&jobName, "job", "", "name of the job to run")
}

func runStep(cmd *cobra.Command, args []string) error {
	var (
		ctx  = cmd.Context()
		s    *workflow.Step
		j    *workflow.Job
		w    *workflow.Workflow
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

	if stepID != "" {
		if jobName != "" {
			w, err := workflow.NewWorkflowFromReader(r)
			if err != nil {
				return err
			}

			j, err = w.GetJob(jobName)
			if err != nil {
				return err
			}
		} else {
			j, err = workflow.NewJobFromReader(r)
			if err != nil {
				return err
			}
		}

		s, err = j.GetStep(stepID)
		if err != nil {
			return err
		}
	} else {
		s, err = workflow.NewStepFromReader(r)
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
		sequence.WithJob(j),
		sequence.WithWorkflow(w),
		sequence.WithRunnerImage(c.Runtime.RunnerImage),
		sequence.WithRepository(flags.FlagWorkDir),
	}
	if c.Verbose {
		opts = append(opts, sequence.WithVerbose)
	}

	return client.RunStep(ctx, s, log.Writer(), opts...)
}
