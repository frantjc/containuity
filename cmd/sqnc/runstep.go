package main

import (
	"fmt"
	"os"

	"github.com/frantjc/sequence"
	"github.com/spf13/cobra"
)

var runStepCmd = &cobra.Command{
	RunE:          runRunStep,
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "step",
}

func init() {
	runStepCmd.Flags().StringVarP(&stepID, "id", "s", "", "ID of the step to run")
	runStepCmd.Flags().StringVarP(&jobName, "job", "j", "", "Name of the job to run")
	runStepCmd.Flags().StringVarP(&runtimeName, "runtime", "", "docker", "Container runtime to use")
}

func runRunStep(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		var (
			ctx       = cmd.Context()
			step      *sequence.Step
			job       *sequence.Job
			path      = args[0]
			file, err = os.Open(path)
		)
		if err != nil {
			return err
		}

		if stepID != "" {
			if jobName != "" {
				ctx = withJob(ctx, encode(file.Name(), jobName, stepID))
				workflow, err := sequence.NewWorkflowFromReader(file)
				if err != nil {
					return err
				}

				job, err = workflow.Job(jobName)
				if err != nil {
					return err
				}
			} else {
				ctx = withJob(ctx, encode(file.Name(), stepID))
				job, err = sequence.NewJobFromReader(file)
				if err != nil {
					return err
				}
			}

			step, err = job.Step(stepID)
			if err != nil {
				return err
			}
		} else {
			ctx = withJob(ctx, encode(file.Name()))
			step, err = sequence.NewStepFromReader(file)
			if err != nil {
				return err
			}
		}

		runtime, err := sequence.GetRuntime(ctx, runtimeName)
		if err != nil {
			return err
		}

		err = runtime.Run(ctx, step)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("not enough arguments")
	}

	return nil
}
