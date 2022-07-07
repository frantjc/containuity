package actions_test

import (
	"testing"

	"github.com/frantjc/sequence/pkg/github/actions"
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
		expected = []string{"::set-output name=var,otherParam=param::value", "::set-output otherParam=param,name=var::value"}
		actual   = command.String()
	)
	assert.Contains(t, expected, actual)
}
