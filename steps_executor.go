package sequence

import (
	"context"
	"os"

	"github.com/frantjc/go-js"
	"github.com/frantjc/sequence/internal/paths"
	"github.com/frantjc/sequence/internal/paths/volumes"
	"github.com/frantjc/sequence/pkg/github/actions"
	"github.com/frantjc/sequence/runtime"
)

type stepsExecutor struct {
	*executor
	steps            []*Step
	preStepWrappers  []*stepWrapper
	mainStepWrappers []*stepWrapper
	postStepWrappers []*stepWrapper
}

func NewStepsExecutor(ctx context.Context, steps []*Step, opts ...ExecutorOpt) (Executor, error) {
	var (
		gc, err = actions.NewContext(paths.GlobalContextOpts()...)
		e       = &stepsExecutor{
			executor: &executor{
				Stdin:         os.Stdin,
				Stdout:        os.Stdout,
				Stderr:        os.Stderr,
				RunnerImage:   DefaultRunnerImage,
				GlobalContext: gc,
			},
			steps: steps,
		}
	)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err := opt(e.executor); err != nil {
			return nil, err
		}
	}

	return e, nil
}

func (e *stepsExecutor) Execute(ctx context.Context) error {
	for _, step := range e.steps {
		if step.IsGitHubAction() {
			action, err := actions.ParseReference(step.Uses)
			if err != nil {
				return err
			}

			actionMetadata, err := e.SetupAction(ctx, action)
			if err != nil {
				return err
			}

			if actionMetadata.IsComposite() {
				steps, err := NewStepsFromCompositeActionMetadata(actionMetadata, paths.Action)
				if err != nil {
					return err
				}

				re := &stepsExecutor{
					executor: e.executor,
					steps:    steps,
				}

				if err := re.Execute(ctx); err != nil {
					return err
				}

				e.preStepWrappers = append(e.preStepWrappers, re.preStepWrappers...)
				e.mainStepWrappers = append(e.mainStepWrappers, re.mainStepWrappers...)
				e.postStepWrappers = append(js.Reverse(re.postStepWrappers), e.postStepWrappers...)
			} else {
				preStep, mainStep, postStep, err := NewStepsFromNonCompositeMetadata(actionMetadata, paths.Action, step)
				if err != nil {
					return err
				}

				var (
					extraMounts = []*runtime.Mount{
						{
							Source:      volumes.GetActionSource(action),
							Destination: e.GlobalContext.GitHubContext.ActionPath,
							Type:        runtime.MountTypeVolume,
						},
					}
					state = map[string]string{}
				)

				if preStep != nil {
					e.preStepWrappers = append(e.preStepWrappers, &stepWrapper{
						step:        preStep,
						extraMounts: extraMounts,
						state:       state,
						id:          e.ID,
					})
				}

				if mainStep != nil {
					e.mainStepWrappers = append(e.mainStepWrappers, &stepWrapper{
						step:        mainStep,
						extraMounts: extraMounts,
						state:       state,
						id:          e.ID,
					})
				}

				if postStep != nil {
					e.postStepWrappers = append([]*stepWrapper{
						{
							step:        postStep,
							extraMounts: extraMounts,
							state:       state,
							id:          e.ID,
						},
					}, e.postStepWrappers...)
				}
			}
		} else {
			e.mainStepWrappers = append(e.mainStepWrappers, &stepWrapper{
				step: step,
			})
		}
	}

	for _, stepWrapper := range append(
		append(e.preStepWrappers, e.mainStepWrappers...),
		e.postStepWrappers...,
	) {
		swe := &stepWrapperExecutor{
			executor:           e.executor,
			stepWrapper:        stepWrapper,
			stopCommandsTokens: map[string]bool{},
		}

		if err := swe.ExecuteStep(ctx); err != nil {
			return err
		}

		e.executor = swe.executor
	}

	return nil
}
