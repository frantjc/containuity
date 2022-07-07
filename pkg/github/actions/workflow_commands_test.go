package actions_test

import (
	"testing"

	"github.com/frantjc/sequence/pkg/github/actions"
	"github.com/stretchr/testify/assert"
)

func TestParseCommandNoParams(t *testing.T) {
	var (
		command  = "::debug::hello there"
		expected = &actions.WorkflowCommand{
			Command:    "debug",
			Parameters: map[string]string{},
			Value:      "hello there",
		}
		actual, err = actions.ParseStringWorkflowCommand(command)
	)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseCommandOneParam(t *testing.T) {
	var (
		command  = "::save-state name=isPost::true"
		expected = &actions.WorkflowCommand{
			Command: "save-state",
			Parameters: map[string]string{
				"name": "isPost",
			},
			Value: "true",
		}
		actual, err = actions.ParseStringWorkflowCommand(command)
	)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseCommandManyParams(t *testing.T) {
	var (
		command  = "::save-state name=isPost,otherParam=1::true"
		expected = &actions.WorkflowCommand{
			Command: "save-state",
			Parameters: map[string]string{
				"name":       "isPost",
				"otherParam": "1",
			},
			Value: "true",
		}
		actual, err = actions.ParseStringWorkflowCommand(command)
	)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseCommandNoValue(t *testing.T) {
	var (
		command  = "::save-state name=isPost::"
		expected = &actions.WorkflowCommand{
			Command: "save-state",
			Parameters: map[string]string{
				"name": "isPost",
			},
			Value: "",
		}
		actual, err = actions.ParseStringWorkflowCommand(command)
	)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
