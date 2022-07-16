package uses_test

import (
	"testing"

	"github.com/frantjc/sequence/pkg/github/actions/uses"
	"github.com/stretchr/testify/assert"
)

func TestParseAction(t *testing.T) {
	var (
		usesStr            = "actions/checkout@v2"
		expectedOwner      = "actions"
		expectedRepository = "checkout"
		expectedPath       = ""
		expectedVersion    = "v2"
		actual, err        = uses.Parse(usesStr)
	)
	assert.Nil(t, err)

	assert.Equal(t, expectedOwner, actual.Owner)
	assert.Equal(t, expectedRepository, actual.Repository)
	assert.Equal(t, expectedPath, actual.Path)
	assert.Equal(t, expectedVersion, actual.Version)
	assert.Equal(t, usesStr, actual.String())
}

func TestParseActionWithPath(t *testing.T) {
	var (
		usesStr            = "frantjc/sequence/testdata@main"
		expectedOwner      = "frantjc"
		expectedRepository = "sequence"
		expectedPath       = "testdata"
		expectedVersion    = "main"
		actual, err        = uses.Parse(usesStr)
	)
	assert.Nil(t, err)

	assert.Equal(t, expectedOwner, actual.Owner)
	assert.Equal(t, expectedRepository, actual.Repository)
	assert.Equal(t, expectedPath, actual.Path)
	assert.Equal(t, expectedVersion, actual.Version)
	assert.Equal(t, usesStr, actual.String())
}
