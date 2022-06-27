package e2e_test

import (
	"bytes"
	"context"
	"os"
	"path"
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

func TestStepExecutorCheckout(t *testing.T) {
	var (
		ctx               = context.TODO()
		imagesPulled      = []runtime.Image{}
		containersCreated = []runtime.Container{}
		volumesCreated    = []runtime.Volume{}
		step, err         = sequence.NewStepFromReader(
			bytes.NewReader(testdata.CheckoutStep),
		)
	)
	assert.Nil(t, err)
	assert.NotNil(t, step)

	se, err := NewTestStepsExecutor(
		t, []*sequence.Step{step},
		sequence.OnImagePull(func(i runtime.Image) {
			imagesPulled = append(imagesPulled, i)
		}),
		sequence.OnContainerCreate(func(c runtime.Container) {
			containersCreated = append(containersCreated, c)
		}),
		sequence.OnVolumeCreate(func(v runtime.Volume) {
			volumesCreated = append(volumesCreated, v)
		}),
	)
	assert.Nil(t, err)
	assert.NotNil(t, se)

	err = se.Execute(ctx)
	assert.Nil(t, err)
	assert.Greater(t, len(imagesPulled), 0)
	assert.True(t, js.Some(imagesPulled, func(i runtime.Image, _ int, _ []runtime.Image) bool {
		return i.GetRef() == sequence.ImageNode12.GetRef()
	}))
	assert.Greater(t, len(containersCreated), 0)
	assert.Greater(t, len(volumesCreated), 0)

	action, err := actions.ParseReference(step.Uses)
	assert.Nil(t, err)
	assert.True(t, js.Some(volumesCreated, func(v runtime.Volume, _ int, _ []runtime.Volume) bool {
		return volumes.GetActionSource(action) == v.GetSource()
	}))
}

func NewTestStepsExecutor(t *testing.T, steps []*sequence.Step, opts ...sequence.ExecutorOpt) (*sequence.StepsExecutor, error) {
	var (
		ctx         = context.TODO()
		githubToken = os.Getenv("SQNC_GITHUB_TOKEN")
		wd, err     = os.Getwd()
	)
	assert.Nil(t, err)
	assert.NotEmpty(t, githubToken)

	rt, err := docker.NewRuntime(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, rt)

	// $GOPATH/src/github.com/frantjc/sequence/internal/e2e
	// => $GOPATH/src/github.com/frantjc/sequence
	rp := path.Dir(path.Dir(wd))

	gc, err := actions.NewContextFromPath(ctx, rp, actions.WithToken(githubToken))
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
