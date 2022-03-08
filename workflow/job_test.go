package workflow_test

import (
	"bytes"
	"testing"

	"github.com/frantjc/sequence/testdata"
	"github.com/frantjc/sequence/workflow"
	"github.com/stretchr/testify/assert"
)

func TestNewJobFromReader(t *testing.T) {
	var (
		// expected = &sequence.Job{}
		_, err = workflow.NewJobFromReader(bytes.NewReader(testdata.CheckoutTestJob))
	)
	assert.Nil(t, err)

	// assert.Equal(t, expected, actual)
}
