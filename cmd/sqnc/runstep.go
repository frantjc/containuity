package main

import (
	"os"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/pkg/orchestrator"
	"github.com/frantjc/sequence/pkg/runtime"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var runStepCmd = &cobra.Command{
	RunE: runRunStep,
	Use:  "step",
	Args: cobra.MinimumNArgs(1),
}

func init() {
	runStepCmd.Flags().StringVarP(&stepID, "id", "s", "", "ID of the step to run")
	runStepCmd.Flags().StringVarP(&jobName, "job", "j", "", "Name of the job to run")
	runStepCmd.Flags().StringVarP(&runtimeName, "runtime", "", "docker", "Container runtime to use")
}

func runRunStep(cmd *cobra.Command, args []string) error {
	var (
		ctx       = cmd.Context()
		step      *sequence.Step
		job       *sequence.Job
		path      = args[0]
		file, err = os.Open(path)
	)
	if err != nil {
		log.Debug().Err(err).Msgf("opening file failed %s", path)
		return err
	}

	if stepID != "" {
		log.Debug().Msg("--id non-empty, must be job or workflow")
		if jobName != "" {
			log.Debug().Msg("--job non-empty, must be workflow")
			workflow, err := sequence.NewWorkflowFromReader(file)
			if err != nil {
				log.Debug().Err(err).Msgf("parsing workflow failed %s", path)
				return err
			}

			job, err = workflow.GetJob(jobName)
			if err != nil {
				log.Debug().Err(err).Msgf("getting job failed %s | %s", path, jobName)
				return err
			}
		} else {
			log.Debug().Msg("--job empty, must be job")
			job, err = sequence.NewJobFromReader(file)
			if err != nil {
				log.Debug().Err(err).Msgf("parsing job failed %s", path)
				return err
			}
		}

		step, err = job.GetStep(stepID)
		if err != nil {
			log.Debug().Err(err).Msgf("getting step failed %s | %s", path, stepID)
			return err
		}
	} else {
		log.Debug().Msg("--step empty, must be step")
		step, err = sequence.NewStepFromReader(file)
		if err != nil {
			log.Debug().Err(err).Msgf("parsing step failed %s", path)
			return err
		}
	}

	r, err := runtime.Get(ctx, runtimeName)
	if err != nil {
		log.Debug().Err(err).Msgf("getting runtime %s", runtimeName)
		return err
	}

	return orchestrator.RunStep(ctx, r, step)
}
