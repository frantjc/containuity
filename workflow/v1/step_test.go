package workflowv1_test

import (
	"bytes"
	"testing"

	"github.com/frantjc/sequence/testdata"
	workflowv1 "github.com/frantjc/sequence/workflow/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewStepFromReader(t *testing.T) {
	var (
		// expected = &sequence.Step{}
		actual, err = workflowv1.NewStepFromReader(bytes.NewReader(testdata.EnvStep))
	)
	assert.Nil(t, err)
	assert.False(t, actual.IsGitHubAction())
	assert.False(t, actual.IsStdoutResponse())

	// assert.Equal(t, expected, actual)
}

func TestUsesNewStepFromReader(t *testing.T) {
	var (
		// expected = &sequence.Step{}
		actual, err = workflowv1.NewStepFromReader(bytes.NewReader(testdata.CheckoutStep))
	)
	assert.Nil(t, err)
	assert.True(t, actual.IsGitHubAction())
	assert.True(t, actual.IsStdoutResponse())

	// assert.Equal(t, expected, actual)
}
