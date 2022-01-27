package main

import (
	"io"
	"os"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/orchestrator"
	"github.com/frantjc/sequence/runtime"
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
	runStepCmd.Flags().StringVarP(&jobName, "job", "j", "", "name of the job to run")
}

func runRunStep(cmd *cobra.Command, args []string) error {
	var (
		ctx  = cmd.Context()
		step *sequence.Step
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
			log.Debug().Err(err).Msgf("opening file failed %s", path)
			return err
		}
	}

	if stepID != "" {
		if jobName != "" {
			workflow, err := sequence.NewWorkflowFromReader(r)
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
			job, err = sequence.NewJobFromReader(r)
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
		step, err = sequence.NewStepFromReader(r)
		if err != nil {
			log.Debug().Err(err).Msgf("parsing step failed %s", path)
			return err
		}
	}

	rt, err := runtime.Get(ctx, runtimeName)
	if err != nil {
		log.Debug().Err(err).Msgf("getting runtime %s", runtimeName)
		return err
	}

	return orchestrator.RunStep(ctx, rt, step, orchestrator.WithGitHubToken(gitHubToken))
}
