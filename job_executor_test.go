package sequence_test

import (
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

func TestJobExecutorCheckout(t *testing.T) {
	for _, rt := range NewTestRuntimes(t) {
		JobExecutorTest(
			t, rt,
			&sequence.Job{
				Steps: []*sequence.Step{
					{
						Uses: "actions/checkout@v3",
					},
				},
			},
			sequence.OnImagePull(func(event *sequence.Event[runtime.Image]) {
				assert.Contains(t, []string{sequence.ImageNode16.GetRef(), sequence.ImageNode12.GetRef()}, event.Type.GetRef())
			}),
		)
	}
}

func TestJobExecutorContainerImage(t *testing.T) {
	for _, rt := range NewTestRuntimes(t) {
		JobExecutorTest(
			t, rt,
			&sequence.Job{
				Container: &sequence.Job_Container{
					Image: golang118Ref,
				},
				Steps: []*sequence.Step{
					{
						Run: "go version",
					},
				},
			},
			sequence.OnImagePull(func(event *sequence.Event[runtime.Image]) {
				assert.Contains(t, []string{sequence.ImageNode16.GetRef(), sequence.ImageNode12.GetRef(), golang118Ref}, event.Type.GetRef())
			}),
		)
	}
}

func TestJobExecutorEnv(t *testing.T) {
	value := "general kenobi"
	for _, rt := range NewTestRuntimes(t) {
		JobExecutorTest(
			t, rt,
			&sequence.Job{
				Env: map[string]string{
					"HELLO_THERE": value,
				},
				Steps: []*sequence.Step{
					{
						Run: "echo ::debug::$HELLO_THERE",
					},
				},
			},
			sequence.OnWorkflowCommand(func(event *sequence.Event[*actions.WorkflowCommand]) {
				assert.Equal(t, value, event.Type.Value)
			}),
		)
	}
}

func JobExecutorTest(t *testing.T, rt runtime.Runtime, job *sequence.Job, opts ...sequence.ExecutorOpt) {
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

	je, err := NewTestJobExecutor(
		t, rt, job,
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
	assert.NotNil(t, je)
	assert.Nil(t, err)

	assert.Nil(t, je.Execute(ctx))
	assert.Greater(t, len(imagesPulled), 0)
	assert.Greater(t, len(containersCreated), 0)
	assert.Greater(t, len(volumesCreated), 0)

	for _, step := range job.Steps {
		if step.IsGitHubAction() {
			action, err := uses.Parse(step.Uses)
			assert.Nil(t, err)
			assert.True(t, js.Some(volumesCreated, func(v runtime.Volume, _ int, _ []runtime.Volume) bool {
				return volumes.GetActionSource(action) == v.GetSource()
			}))
		}
	}
}

func NewTestJobExecutor(t *testing.T, rt runtime.Runtime, job *sequence.Job, opts ...sequence.ExecutorOpt) (sequence.Executor, error) {
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

	je, err := sequence.NewJobExecutor(
		ctx, job,
		append([]sequence.ExecutorOpt{
			sequence.WithRuntime(rt),
			sequence.WithGlobalContext(gc),
		}, opts...)...,
	)
	assert.NotNil(t, je)
	assert.Nil(t, err)

	return je, err
}
