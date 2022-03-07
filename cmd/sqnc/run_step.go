package main

import (
	"fmt"
	"io"
	"os"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/conf"
	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/sio"
	"github.com/frantjc/sequence/workflow"
	"github.com/spf13/cobra"
)

var runStepCmd = &cobra.Command{
	RunE: runStep,
	Use:  "step",
	Args: cobra.ExactArgs(1),
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

	if stepID != "" {
		if jobName != "" {
			workflow, err := workflow.NewWorkflowFromReader(r)
			if err != nil {
				return err
			}

			j, err = workflow.GetJob(jobName)
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

	c, err := conf.Get()
	if err != nil {
		return err
	}

	client, err := sequence.New(ctx, c.Address)
	if err != nil {
		return err
	}

	return client.RunStep(ctx, s, sio.NewPrefixedWriter(fmt.Sprintf("%s|%s ", log.ColorInfo, log.ColorNone), log.Writer()))
}
