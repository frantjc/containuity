package sequence

import (
	"context"
	"os"

	"github.com/frantjc/go-js"
	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/runtime"
)

type StepsExecutor struct {
	Executor
	Steps            []*Step
	preStepWrappers  []*StepWrapper
	mainStepWrappers []*StepWrapper
	postStepWrappers []*StepWrapper
}

func NewStepsExecutor(ctx context.Context, steps []*Step, opts ...ExecutorOpt) (*StepsExecutor, error) {
	var (
		gc, err = actions.NewContext(defaultGlobalContextOpts()...)
		e       = &StepsExecutor{
			Executor: Executor{
				Stdin:         os.Stdin,
				Stdout:        os.Stdout,
				Stderr:        os.Stderr,
				RunnerImage:   DefaultRunnerImage,
				GlobalContext: gc,
			},
			Steps: steps,
		}
	)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err := opt(&e.Executor); err != nil {
			return nil, err
		}
	}

	return e, nil
}

func (e *StepsExecutor) Execute(ctx context.Context) error {
	for _, step := range e.Steps {
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
				steps, err := NewStepsFromCompositeActionMetadata(actionMetadata, actionPath)
				if err != nil {
					return err
				}

				re := &StepsExecutor{
					Executor: e.Executor,
					Steps:    steps,
				}

				if err := re.Execute(ctx); err != nil {
					return err
				}

				e.preStepWrappers = append(e.preStepWrappers, re.preStepWrappers...)
				e.mainStepWrappers = append(e.mainStepWrappers, re.mainStepWrappers...)
				e.postStepWrappers = append(js.Reverse(re.postStepWrappers), e.postStepWrappers...)
			} else {
				preStep, mainStep, postStep, err := NewStepsFromNonCompositeMetadata(actionMetadata, actionPath, step)
				if err != nil {
					return err
				}

				var (
					extraMounts = []*runtime.Mount{
						{
							Source:      GetActionVolumeName(action),
							Destination: e.GlobalContext.GitHubContext.ActionPath,
							Type:        runtime.MountTypeVolume,
						},
					}
					state = map[string]string{}
				)

				if preStep != nil {
					e.preStepWrappers = append(e.preStepWrappers, &StepWrapper{
						Step:        preStep,
						ExtraMounts: extraMounts,
						State:       state,
					})
				}

				if mainStep != nil {
					e.mainStepWrappers = append(e.mainStepWrappers, &StepWrapper{
						Step:        mainStep,
						ExtraMounts: extraMounts,
						State:       state,
					})
				}

				if postStep != nil {
					e.postStepWrappers = append([]*StepWrapper{
						{
							Step:        postStep,
							ExtraMounts: extraMounts,
							State:       state,
						},
					}, e.postStepWrappers...)
				}
			}
		} else {
			e.mainStepWrappers = append(e.mainStepWrappers, &StepWrapper{
				Step: step,
			})
		}
	}

	for _, stepWrapper := range append(
		append(e.preStepWrappers, e.mainStepWrappers...),
		e.postStepWrappers...,
	) {
		se := &StepExecutor{
			Executor:    e.Executor,
			StepWrapper: stepWrapper,
		}

		if err := se.ExecuteStep(ctx); err != nil {
			return err
		}
	}

	return nil
}
