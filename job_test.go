package sequence_test

import (
	"bytes"
	"testing"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJobFromReader(t *testing.T) {
	var (
		// expected = &sequence.Job{}
		_, err = sequence.NewJobFromReader(bytes.NewReader(testdata.Job))
	)
	assert.Nil(t, err)

	// assert.Equal(t, expected, actual)
}
