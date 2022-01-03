package sequence_test

import (
	"bytes"
	"testing"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewWorkflowFromString(t *testing.T) {
	var (
		// expected = &sequence.Workflow{}
		_, err = sequence.NewWorkflowFromString(string(testdata.Workflow))
	)
	assert.Nil(t, err)

	// assert.Equal(t, expected, actual)
}


func TestNewWorkflowFromBytes(t *testing.T) {
	var (
		// expected = &sequence.Workflow{}
		_, err = sequence.NewWorkflowFromBytes(testdata.Workflow)
	)
	assert.Nil(t, err)

	// assert.Equal(t, expected, actual)
}

func TestNewWorkflowFromReader(t *testing.T) {
	var (
		// expected = &sequence.Workflow{}
		_, err = sequence.NewWorkflowFromReader(bytes.NewReader(testdata.Workflow))
	)
	assert.Nil(t, err)

	// assert.Equal(t, expected, actual)
}
