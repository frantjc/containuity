package actions_test

import (
	"testing"

	"github.com/frantjc/sequence/internal/actions"
	"github.com/stretchr/testify/assert"
)

func TestParseAction(t *testing.T) {
	var (
		action   = "actions/checkout@v2"
		expected = &actions.Uses{
			Owner:      "actions",
			Repository: "checkout",
			Path:       "",
			Version:    "v2",
		}
		actual, err = actions.Parse(action)
	)
	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
	assert.Equal(t, "actions/checkout", actual.Repo())
	assert.Equal(t, action, actual.String())
}

func TestParseActionWithPath(t *testing.T) {
	var (
		action   = "frantjc/sequence/internal/testdata@main"
		expected = &actions.Uses{
			Owner:      "frantjc",
			Repository: "sequence",
			Path:       "internal/testdata",
			Version:    "main",
		}
		actual, err = actions.Parse(action)
	)
	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
	assert.Equal(t, "frantjc/sequence", actual.Repo())
	assert.Equal(t, action, actual.String())
}
