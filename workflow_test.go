package sequence_test

import (
	"bytes"
	"testing"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewWorkflowFromReader(t *testing.T) {
	// expected = &sequence.Workflow{}
	_, err := sequence.NewWorkflowFromReader(bytes.NewReader(testdata.CheckoutTestBuildWorkflow))
	assert.Nil(t, err)

	// assert.Equal(t, expected, actual)
}
