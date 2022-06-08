package workflowv1_test

import (
	"bytes"
	"testing"

	"github.com/frantjc/sequence/testdata"
	workflowv1 "github.com/frantjc/sequence/workflow/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewJobFromReader(t *testing.T) {
	var (
		// expected = &sequence.Job{}
		_, err = workflowv1.NewJobFromReader(bytes.NewReader(testdata.CheckoutTestJob))
	)
	assert.Nil(t, err)

	// assert.Equal(t, expected, actual)
}
