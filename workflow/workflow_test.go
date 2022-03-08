package workflow_test

import (
	"bytes"
	"testing"

	"github.com/frantjc/sequence/testdata"
	"github.com/frantjc/sequence/workflow"
	"github.com/stretchr/testify/assert"
)

func TestNewWorkflowFromReader(t *testing.T) {
	var (
		// expected = &sequence.Workflow{}
		_, err = workflow.NewWorkflowFromReader(bytes.NewReader(testdata.CheckoutTestBuildWorkflow))
	)
	assert.Nil(t, err)

	// assert.Equal(t, expected, actual)
}
