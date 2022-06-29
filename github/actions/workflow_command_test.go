package actions_test

import (
	"testing"

	"github.com/frantjc/sequence/github/actions"
	"github.com/stretchr/testify/assert"
)

func TestWorkflowCommandToString(t *testing.T) {
	var (
		command = &actions.WorkflowCommand{
			Command: "set-output",
			Parameters: map[string]string{
				"name":       "var",
				"otherParam": "param",
			},
			Value: "value",
		}
		expected1 = "::set-output name=var,otherParam=param::value"
		expected2 = "::set-output otherParam=param,name=var::value"
		actual    = command.String()
	)
	assert.True(t, actual == expected1 || actual == expected2)
}
