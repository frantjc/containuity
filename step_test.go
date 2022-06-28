package sequence_test

import (
	"bytes"
	"testing"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewStepFromReader(t *testing.T) {
	var (
		actual, err = sequence.NewStepFromReader(bytes.NewReader(testdata.EnvStep))
	)
	assert.Nil(t, err)
	assert.False(t, actual.IsGitHubAction())
}

func TestUsesNewStepFromReader(t *testing.T) {
	var (
		actual, err = sequence.NewStepFromReader(bytes.NewReader(testdata.CheckoutStep))
	)
	assert.Nil(t, err)
	assert.True(t, actual.IsGitHubAction())
}
