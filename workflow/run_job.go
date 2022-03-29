package workflow

import (
	"context"
	"fmt"
	"os"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/runtime"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func RunJob(ctx context.Context, r runtime.Runtime, j *Job, opts ...RunOpt) error {
	ro, err := newRunOpts(append(opts, WithJob(j))...)
	if err != nil {
		return err
	}

	ctxopts := []actions.CtxOpt{
		actions.WithToken(ro.githubToken),
		actions.WithSecrets(ro.secrets),
		actions.WithWorkdir(containerWorkdir),
		actions.WithJobName(ro.jobName),
	}
	ro.gctx, err = actions.NewContextFromPath(ctx, ro.repository, ctxopts...)
	if err != nil {
		return err
	}

	ro.logout.Infof("[%sSQNC%s] running job '%s'", log.ColorInfo, log.ColorNone, j.Name)
	return runJob(ctx, r, j, ro)
}

func runJob(ctx context.Context, r runtime.Runtime, j *Job, ro *runOpts) error {
	var (
		jobID          = getID(ro)
		jobWorkdir     = getHostWorkdir(jobID, ro)
		hostGitHubEnv  = getHostGitHubEnvFilepath(jobWorkdir)
		hostGitHubPath = getHostGitHubPathFilepath(jobWorkdir)
	)
	os.Remove(hostGitHubEnv)
	os.Remove(hostGitHubPath)
	defer os.Remove(hostGitHubEnv)
	defer os.Remove(hostGitHubPath)

	var (
		preSteps  []*Step
		mainSteps []*Step
		postSteps []*Step
	)
	for _, step := range j.Steps {
		if step.Uses != "" {
			expandedStep, err := expandStep(step.Canonical(), ro.gctx)
			if err != nil {
				return err
			}

			if metadata, reference, ro, err := runStepSetup(ctx, r, expandedStep, ro); err != nil {
				return err
			} else if metadata != nil && reference != nil {
				ro.specOpts = append(ro.specOpts, runtime.WithMounts(
					specs.Mount{
						Source:      getHostActionPath(reference, ro),
						Destination: ro.gctx.GitHubContext.ActionPath,
						Type:        runtime.MountTypeBind,
					},
				))
				with := withFromInputs(metadata.Inputs)
				if preStep, err := newPreStepFromMetadataWith(metadata, ro.gctx.GitHubContext.ActionPath, with); err != nil {
					return err
				} else if preStep != nil {
					ro.logout.Debugf("[%sSQNC:DBG%s] adding action pre step '%s'", log.ColorDebug, log.ColorNone, expandedStep.GetID())
					if ro.actionImage != "" {
						preStep.Image = ro.actionImage
					}
					if expandedStep.ID != "" {
						preStep.ID = fmt.Sprintf("Pre %s", expandedStep.ID)
					}
					if expandedStep.Name != "" {
						preStep.Name = fmt.Sprintf("Pre %s", expandedStep.Name)
					}
					preSteps = append(preSteps, preStep)
				}

				if mainStep, err := newMainStepFromMetadataWith(metadata, ro.gctx.GitHubContext.ActionPath, with); err != nil {
					return err
				} else if mainStep != nil {
					ro.logout.Debugf("[%sSQNC:DBG%s] adding action main step '%s'", log.ColorDebug, log.ColorNone, expandedStep.GetID())
					if ro.actionImage != "" {
						mainStep.Image = ro.actionImage
					}
					mainStep.ID = expandedStep.ID
					mainStep.Name = expandedStep.Name
					mainSteps = append(mainSteps, mainStep)
				}

				if postStep, err := newPostStepFromMetadataWith(metadata, ro.gctx.GitHubContext.ActionPath, with); err != nil {
					return err
				} else if postStep != nil {
					ro.logout.Debugf("[%sSQNC:DBG%s] adding action post step '%s'", log.ColorDebug, log.ColorNone, expandedStep.GetID())
					if ro.actionImage != "" {
						postStep.Image = ro.actionImage
					}
					if expandedStep.ID != "" {
						postStep.ID = fmt.Sprintf("Post %s", expandedStep.ID)
					}
					if expandedStep.Name != "" {
						postStep.Name = fmt.Sprintf("Post %s", expandedStep.Name)
					}
					postSteps = append(postSteps, postStep)
				}
			}
		} else {
			ro.logout.Debugf("[%sSQNC:DBG%s] adding raw step '%s'", log.ColorDebug, log.ColorNone, step.GetID())
			mainSteps = append(mainSteps, step)
		}
	}

	for _, step := range append(append(preSteps, mainSteps...), postSteps...) {
		expandedStep, err := expandStep(step.Canonical(), ro.gctx)
		if err != nil {
			return err
		}

		ro.logout.Infof("[%sSQNC%s] running step '%s'", log.ColorInfo, log.ColorNone, expandedStep.GetID())
		if _, ro, err = runStep(ctx, r, expandedStep, ro); err != nil {
			return err
		}
	}

	return nil
}
