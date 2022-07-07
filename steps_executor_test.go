package sequence_test

import (
	"os"
	"testing"

	"github.com/frantjc/go-js"
	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/paths/volumes"
	"github.com/frantjc/sequence/pkg/github/actions"
	"github.com/frantjc/sequence/runtime"
	"github.com/stretchr/testify/assert"
)

func TestStepsExecutor(t *testing.T) {
	for _, r := range NewTestRuntimes(t) {
		for _, f := range []RuntimeTest{
			StepsExecutorCheckoutSetupGoTest,
			StepsExecutorDefaultImageTest,
			StepsExecutorImageTest,
			StepsExecutoGitHubPathTest,
			StepsExecutoGitHubEnvTest,
			StepsExecutorStopCommandsTest,
			StepsExecutorSetOutputTest,
			PruneTest,
		} {
			f(t, r)
		}
	}
}

func StepsExecutorCheckoutSetupGoTest(t *testing.T, rt runtime.Runtime) {
	StepsExecutorTest(t, rt, []*sequence.Step{
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
			// hilariously, "recursively" run some of sequence's test :)
			Run: "go test ./internal/...",
		},
	})
}

func StepsExecutorDefaultImageTest(t *testing.T, rt runtime.Runtime) {
	StepsExecutorTest(
		t, rt,
		[]*sequence.Step{
			{
				Run: "echo test",
			},
		},
		sequence.WithRunnerImage(alpineImg),
		sequence.OnImagePull(func(i runtime.Image) {
			assert.Equal(t, alpineRef, i.GetRef())
		}),
	)
}

func StepsExecutorImageTest(t *testing.T, rt runtime.Runtime) {
	StepsExecutorTest(
		t, rt,
		[]*sequence.Step{
			{
				Image: alpineRef,
				Run:   "echo test",
			},
		},
		sequence.OnImagePull(func(i runtime.Image) {
			assert.Equal(t, alpineRef, i.GetRef())
		}),
	)
}

func StepsExecutoGitHubPathTest(t *testing.T, rt runtime.Runtime) {
	StepsExecutorTest(
		t, rt,
		[]*sequence.Step{
			{
				Run: "echo /.bin >> $GITHUB_PATH",
			},
			{
				Run: "echo \"::debug::$PATH\"",
			},
		},
		sequence.OnWorkflowCommand(func(wc *actions.WorkflowCommand) {
			assert.Contains(t, wc.Value, "/.bin")
		}),
	)
}

func StepsExecutoGitHubEnvTest(t *testing.T, rt runtime.Runtime) {
	StepsExecutorTest(
		t, rt,
		[]*sequence.Step{
			{
				Run: "echo HELLO_THERE=generalkenobi >> $GITHUB_ENV",
			},
			{
				Run: "echo \"::debug::$HELLO_THERE\"",
			},
		},
		sequence.OnWorkflowCommand(func(wc *actions.WorkflowCommand) {
			assert.Equal(t, wc.Value, "generalkenobi")
		}),
	)
}

func StepsExecutorStopCommandsTest(t *testing.T, rt runtime.Runtime) {
	debugCount := 0

	StepsExecutorTest(
		t, rt,
		[]*sequence.Step{
			{
				Run: `
				echo '::debug::test1'
				echo '::stop-commands::token'
				echo '::debug::test2'
				echo '::token::'
				echo '::debug::test3'
				`,
			},
		},
		sequence.OnWorkflowCommand(func(wc *actions.WorkflowCommand) {
			switch wc.Command {
			case actions.CommandStopCommands:
				assert.Equal(t, "token", wc.Value)
			case actions.CommandDebug:
				debugCount++
			default:
				assert.Equal(t, "token", wc.Command)
			}
		}),
	)

	assert.Equal(t, debugCount, 2)
}

func StepsExecutorSetOutputTest(t *testing.T, rt runtime.Runtime) {
	StepsExecutorTest(
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

func StepsExecutorTest(t *testing.T, rt runtime.Runtime, steps []*sequence.Step, opts ...sequence.ExecutorOpt) {
	var (
		imagesPulled           = []runtime.Image{}
		containersCreated      = []runtime.Container{}
		volumesCreated         = []runtime.Volume{}
		workflowCommandsIssued = []*actions.WorkflowCommand{}
	)

	stdout, err := os.CreateTemp("", "")
	assert.NotNil(t, stdout)
	assert.Nil(t, err)
	defer os.Remove(stdout.Name())

	stderr := stdout
	assert.NotNil(t, stderr)

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
			sequence.WithStreams(os.Stdin, stdout, stderr),
		)...,
	)
	assert.NotNil(t, se)
	assert.Nil(t, err)

	assert.Nil(t, se.Execute(ctx))
	assert.Greater(t, len(imagesPulled), 0)
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
}

func NewTestStepsExecutor(t *testing.T, rt runtime.Runtime, steps []*sequence.Step, opts ...sequence.ExecutorOpt) (sequence.Executor, error) {
	var (
		githubToken = os.Getenv("GITHUB_TOKEN")
		// all tests in this suite are ran against
		// https://github.com/frantjc/sequence
		wd, err = os.Getwd()
	)
	assert.Nil(t, err)
	if !assert.NotEmpty(t, githubToken) {
		assert.FailNow(t, "GITHUB_TOKEN must be set")
	}

	gc, err := actions.NewContextFromPath(ctx, wd, actions.WithToken(githubToken))
	assert.NotNil(t, gc)
	assert.Nil(t, err)

	se, err := sequence.NewStepsExecutor(
		ctx, steps,
		append([]sequence.ExecutorOpt{
			sequence.WithRuntime(rt),
			sequence.WithGlobalContext(gc),
		}, opts...)...,
	)
	assert.NotNil(t, se)
	assert.Nil(t, err)

	return se, err
}
