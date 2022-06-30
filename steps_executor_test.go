package sequence_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/frantjc/go-js"
	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/internal/paths/volumes"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/runtime/docker"
	"github.com/frantjc/sequence/testdata"
	"github.com/stretchr/testify/assert"
)

var (
	ctx = context.TODO()
)

func TestDockerRuntime(t *testing.T) {
	rt, err := docker.NewRuntime(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, rt)

	StepExecutorCheckoutSetupGoTest(t, rt)
}

func StepExecutorCheckoutSetupGoTest(t *testing.T, rt runtime.Runtime) {
	checkoutStep, err := sequence.NewStepFromReader(
		bytes.NewReader(testdata.CheckoutStep),
	)
	assert.Nil(t, err)
	assert.NotNil(t, checkoutStep)

	StepExecutorTest(t, rt, []*sequence.Step{
		{
			Uses: "actions/checkout@v2",
		},
		{
			Uses: "actions/setup-go@v2",
			With: map[string]string{
				"go-version": "1.18",
			},
		},
		{
			// hilariously, recursively run sequence's test :)
			Run: "go test ./github/...",
		},
	})
}

type image string

func (i image) GetRef() string {
	return string(i)
}

const ref = "docker.io/library/alpine"

func StepDefaultImageTest(t *testing.T, rt runtime.Runtime) {
	var (
		img = image(ref)
	)
	StepExecutorTest(
		t, rt,
		[]*sequence.Step{
			{
				Run: "ls",
			},
		},
		sequence.WithRunnerImage(img),
		sequence.OnImagePull(func(i runtime.Image) {
			assert.Equal(t, img.GetRef(), i.GetRef())
		}),
	)
}

func StepImageTest(t *testing.T, rt runtime.Runtime) {
	StepExecutorTest(
		t, rt,
		[]*sequence.Step{
			{
				Image: ref,
				Run:   "ls",
			},
		},
		sequence.OnImagePull(func(i runtime.Image) {
			assert.Equal(t, ref, i.GetRef())
		}),
	)
}

func StepStopCommandsTest(t *testing.T, rt runtime.Runtime) {
	var (
		count = 0
	)
	StepExecutorTest(
		t, rt,
		[]*sequence.Step{
			{
				Run: `
				echo ::debug::test1
				echo ::stop-commands::token
				echo ::debug::test2
				echo ::token::
				echo ::debug::test3
				`,
			},
		},
		sequence.OnWorkflowCommand(func(wc *actions.WorkflowCommand) {
			switch wc.Command {
			case actions.CommandStopCommands:
				assert.Equal(t, "token", wc.Value)
			case actions.CommandDebug:
				count++
			default:
				assert.Equal(t, "token", wc.Command)
			}
		}),
	)
	assert.Equal(t, count, 2)
}

func StepSetOutputTest(t *testing.T, rt runtime.Runtime) {
	StepExecutorTest(
		t, rt,
		[]*sequence.Step{
			{
				Id:  "test",
				Run: "echo '::set-output name=hellothere::general kenobi'",
			},
			{
				Run: "echo '::notice::${{ steps.test.outputs.hellothere }}'",
			},
		},
		sequence.OnWorkflowCommand(func(wc *actions.WorkflowCommand) {
			switch wc.Command {
			case actions.CommandSetOutput:
				assert.Equal(t, "hellothere", wc.Parameters["name"])
				assert.Equal(t, "general kenobi", wc.Value)
			case actions.CommandNotice:
				assert.Equal(t, "general kenobi", wc.Value)
			}
		}),
	)
}

func StepExecutorTest(t *testing.T, rt runtime.Runtime, steps []*sequence.Step, opts ...sequence.ExecutorOpt) {
	var (
		imagesPulled           = []runtime.Image{}
		containersCreated      = []runtime.Container{}
		volumesCreated         = []runtime.Volume{}
		workflowCommandsIssued = []*actions.WorkflowCommand{}
	)

	se, err := NewTestStepsExecutor(
		t, rt, steps,
		append(
			opts,
			sequence.OnImagePull(func(i runtime.Image) {
				imagesPulled = append(imagesPulled, i)
			}),
			sequence.OnContainerCreate(func(c runtime.Container) {
				containersCreated = append(containersCreated, c)
			}),
			sequence.OnVolumeCreate(func(v runtime.Volume) {
				volumesCreated = append(volumesCreated, v)
			}),
			sequence.OnWorkflowCommand(func(wc *actions.WorkflowCommand) {
				workflowCommandsIssued = append(workflowCommandsIssued, wc)
			}),
		)...,
	)
	assert.Nil(t, err)
	assert.NotNil(t, se)

	assert.Nil(t, se.Execute(ctx))
	assert.Greater(t, len(imagesPulled), 0)
	assert.True(t, js.Some(imagesPulled, func(i runtime.Image, _ int, _ []runtime.Image) bool {
		return i.GetRef() == sequence.ImageNode12.GetRef()
	}))
	assert.Greater(t, len(containersCreated), 0)
	assert.Greater(t, len(volumesCreated), 0)

	for _, step := range steps {
		if step.IsGitHubAction() {
			action, err := actions.ParseReference(step.Uses)
			assert.Nil(t, err)
			assert.True(t, js.Some(volumesCreated, func(v runtime.Volume, _ int, _ []runtime.Volume) bool {
				return volumes.GetActionSource(action) == v.GetSource()
			}))
		}
	}

	PruneRuntimeTest(t, rt)
}

func PruneRuntimeTest(t *testing.T, rt runtime.Runtime) {
	assert.Nil(t, rt.PruneContainers(ctx))
	assert.Nil(t, rt.PruneVolumes(ctx))
	assert.Nil(t, rt.PruneImages(ctx))
}

func NewTestStepsExecutor(t *testing.T, rt runtime.Runtime, steps []*sequence.Step, opts ...sequence.ExecutorOpt) (sequence.Executor, error) {
	var (
		githubToken = os.Getenv("SQNC_GITHUB_TOKEN")
		// all tests in this suite are ran against
		// https://github.com/frantjc/sequence
		wd, err = os.Getwd()
	)
	assert.Nil(t, err)
	assert.NotEmpty(t, githubToken)

	gc, err := actions.NewContextFromPath(ctx, wd, actions.WithToken(githubToken))
	assert.Nil(t, err)
	assert.NotNil(t, gc)

	se, err := sequence.NewStepsExecutor(
		ctx, steps,
		append([]sequence.ExecutorOpt{
			sequence.WithRuntime(rt),
			sequence.WithGlobalContext(gc),
		}, opts...)...,
	)
	assert.Nil(t, err)
	assert.NotNil(t, se)

	return se, err
}
