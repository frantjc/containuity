package main

import (
	"fmt"
	"os"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/pkg/runtime"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var runStepCmd = &cobra.Command{
	RunE: runRunStep,
	Use:  "step",
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
				workflow, err := sequence.NewWorkflowFromReader(file)
				if err != nil {
					return err
				}

				job, err = workflow.GetJob(jobName)
				if err != nil {
					return err
				}
			} else {
				job, err = sequence.NewJobFromReader(file)
				if err != nil {
					return err
				}
			}

			step, err = job.GetStep(stepID)
			if err != nil {
				return err
			}
		} else {
			step, err = sequence.NewStepFromReader(file)
			if err != nil {
				return err
			}
		}

		docker, err := runtime.GetRuntime("docker")
		if err != nil {
			return err
		}

		err = docker.Run(ctx, step)
		if err != nil {
			log.Err(err).Msg("run err")
			return err
		}
	} else {
		return fmt.Errorf("not enough arguments")
	}

	return nil
}
