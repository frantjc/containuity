package sequence_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/frantjc/go-js"
	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/paths/volumes"
	"github.com/frantjc/sequence/pkg/github/actions"
	"github.com/frantjc/sequence/pkg/github/actions/uses"
	"github.com/frantjc/sequence/runtime"
	"github.com/stretchr/testify/assert"
)

func TestStepsExecutorCheckoutSetupGo(t *testing.T) {
	for _, rt := range NewTestRuntimes(t) {
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
		})
	}
}

func TestStepsExecutorDefaultImage(t *testing.T) {
	for _, rt := range NewTestRuntimes(t) {
		StepsExecutorTest(
			t, rt,
			[]*sequence.Step{
				{
					Run: "echo test",
				},
			},
			sequence.WithRunnerImage(alpineImg),
			sequence.OnImagePull(func(event *sequence.Event[runtime.Image]) {
				assert.Equal(t, alpineRef, event.Type.GetRef())
			}),
		)
	}
}

func TestStepsExecutorImage(t *testing.T) {
	for _, rt := range NewTestRuntimes(t) {
		StepsExecutorTest(
			t, rt,
			[]*sequence.Step{
				{
					Image: alpineRef,
					Run:   "echo test",
				},
			},
			sequence.OnImagePull(func(event *sequence.Event[runtime.Image]) {
				assert.Equal(t, alpineRef, event.Type.GetRef())
			}),
		)
	}
}

func TestStepsExecutorGitHubPath(t *testing.T) {
	for _, rt := range NewTestRuntimes(t) {
		StepsExecutorTest(
			t, rt,
			[]*sequence.Step{
				{
					Run: fmt.Sprintf("echo /.bin >> $%s", actions.EnvVarPath),
				},
				{
					Run: "echo ::debug::$PATH",
				},
			},
			sequence.OnWorkflowCommand(func(event *sequence.Event[*actions.WorkflowCommand]) {
				assert.Contains(t, event.Type.Value, "/.bin")
			}),
		)
	}
}

func TestStepsExecutorGitHubEnv(t *testing.T) {
	var (
		envVar = "HELLO_THERE"
		value  = "general kenobi"
	)
	for _, rt := range NewTestRuntimes(t) {
		StepsExecutorTest(
			t, rt,
			[]*sequence.Step{
				{
					Run: fmt.Sprintf("echo %s=%s >> $%s", envVar, value, actions.EnvVarEnv),
				},
				{
					Run: fmt.Sprintf("echo ::debug::$%s", envVar),
				},
			},
			sequence.OnWorkflowCommand(func(event *sequence.Event[*actions.WorkflowCommand]) {
				assert.Equal(t, event.Type.Value, value)
			}),
		)
	}
}

func TestStepsExecutorStopCommands(t *testing.T) {
	var (
		stopCommandsToken = "token"
	)
	for _, rt := range NewTestRuntimes(t) {
		var (
			debugCount = 0
		)
		StepsExecutorTest(
			t, rt,
			[]*sequence.Step{
				{
					Run: fmt.Sprintf(`
					echo '::debug::test1'
					echo '::stop-commands::%s'
					echo '::debug::test2'
					echo '::%s::'
					echo '::debug::test3'
					`, stopCommandsToken, stopCommandsToken),
				},
			},
			sequence.OnWorkflowCommand(func(event *sequence.Event[*actions.WorkflowCommand]) {
				switch event.Type.Command {
				case actions.CommandStopCommands:
					assert.Equal(t, stopCommandsToken, event.Type.Value)
				case actions.CommandDebug:
					debugCount++
				default:
					assert.Equal(t, stopCommandsToken, event.Type.Command)
				}
			}),
		)

		assert.Equal(t, debugCount, 2)
	}
}

func TestStepsExecutorSetOutput(t *testing.T) {
	var (
		output = "hellothere"
		value  = "general kenobi"
		stepID = "test"
	)
	for _, rt := range NewTestRuntimes(t) {
		StepsExecutorTest(
			t, rt,
			[]*sequence.Step{
				{
					Id:  stepID,
					Run: fmt.Sprintf("echo ::set-output name=%s::%s", output, value),
				},
				{
					Run: fmt.Sprintf("echo ::notice::${{ steps.%s.outputs.%s }}", stepID, output),
				},
			},
			sequence.OnWorkflowCommand(func(event *sequence.Event[*actions.WorkflowCommand]) {
				switch event.Type.Command {
				case actions.CommandSetOutput:
					assert.Equal(t, output, event.Type.GetName())
					assert.Equal(t, value, event.Type.Value)
				case actions.CommandNotice:
					assert.Equal(t, value, event.Type.Value)
				default:
					assert.True(t, false, "unexpected workflow command", event.Type.String())
				}
			}),
		)
	}
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
			sequence.OnImagePull(func(event *sequence.Event[runtime.Image]) {
				imagesPulled = append(imagesPulled, event.Type)
			}),
			sequence.OnContainerCreate(func(event *sequence.Event[runtime.Container]) {
				containersCreated = append(containersCreated, event.Type)
			}),
			sequence.OnVolumeCreate(func(event *sequence.Event[runtime.Volume]) {
				volumesCreated = append(volumesCreated, event.Type)
			}),
			sequence.OnWorkflowCommand(func(event *sequence.Event[*actions.WorkflowCommand]) {
				workflowCommandsIssued = append(workflowCommandsIssued, event.Type)
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
			action, err := uses.Parse(step.Uses)
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
