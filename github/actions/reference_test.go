package actions_test

import (
	"testing"

	"github.com/frantjc/sequence/github/actions"
	"github.com/stretchr/testify/assert"
)

func TestParseAction(t *testing.T) {
	var (
		actionRef          = "actions/checkout@v2"
		expectedOwner      = "actions"
		expectedRepository = "checkout"
		expectedPath       = ""
		expectedVersion    = "v2"
		actual, err        = actions.ParseReference(actionRef)
	)
	assert.Nil(t, err)

	assert.Equal(t, expectedOwner, actual.Owner())
	assert.Equal(t, expectedRepository, actual.Repository())
	assert.Equal(t, expectedPath, actual.Path())
	assert.Equal(t, expectedVersion, actual.Version())
	assert.Equal(t, actionRef, actual.String())
}

func TestParseActionWithPath(t *testing.T) {
	var (
		actionRef          = "frantjc/sequence/testdata@main"
		expectedOwner      = "frantjc"
		expectedRepository = "sequence"
		expectedPath       = "testdata"
		expectedVersion    = "main"
		actual, err        = actions.ParseReference(actionRef)
	)
	assert.Nil(t, err)

	assert.Equal(t, expectedOwner, actual.Owner())
	assert.Equal(t, expectedRepository, actual.Repository())
	assert.Equal(t, expectedPath, actual.Path())
	assert.Equal(t, expectedVersion, actual.Version())
	assert.Equal(t, actionRef, actual.String())
}
