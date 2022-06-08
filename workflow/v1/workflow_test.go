package workflowv1_test

import (
	"bytes"
	"testing"

	"github.com/frantjc/sequence/testdata"
	workflowv1 "github.com/frantjc/sequence/workflow/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewWorkflowFromReader(t *testing.T) {
	var (
		// expected = &sequence.Workflow{}
		_, err = workflowv1.NewWorkflowFromReader(bytes.NewReader(testdata.CheckoutTestBuildWorkflow))
	)
	assert.Nil(t, err)

	// assert.Equal(t, expected, actual)
}
