package sequence_test

import (
	"bytes"
	"testing"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJobFromString(t *testing.T) {
	var (
		// expected = &sequence.Job{}
		_, err = sequence.NewJobFromString(string(testdata.Job))
	)
	assert.Nil(t, err)

	// assert.Equal(t, expected, actual)
}

func TestNewJobFromBytes(t *testing.T) {
	var (
		// expected = &sequence.Job{}
		_, err = sequence.NewJobFromBytes(testdata.Job)
	)
	assert.Nil(t, err)

	// assert.Equal(t, expected, actual)
}

func TestNewJobFromReader(t *testing.T) {
	var (
		// expected = &sequence.Job{}
		_, err = sequence.NewJobFromReader(bytes.NewReader(testdata.Job))
	)
	assert.Nil(t, err)

	// assert.Equal(t, expected, actual)
}
