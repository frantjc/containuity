package sequence_test

import (
	"bytes"
	"testing"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewStepFromString(t *testing.T) {
	var (
		// expected = &sequence.Step{}
		actual, err = sequence.NewStepFromString(string(testdata.Step))
	)
	assert.Nil(t, err)
	assert.False(t, actual.IsAction())
	assert.False(t, actual.IsStdoutParsable())

	// assert.Equal(t, expected, actual)
}


func TestNewStepFromBytes(t *testing.T) {
	var (
		// expected = &sequence.Step{}
		actual, err = sequence.NewStepFromBytes(testdata.Step)
	)
	assert.Nil(t, err)
	assert.False(t, actual.IsAction())
	assert.False(t, actual.IsStdoutParsable())

	// assert.Equal(t, expected, actual)
}

func TestNewStepFromReader(t *testing.T) {
	var (
		// expected = &sequence.Step{}
		actual, err = sequence.NewStepFromReader(bytes.NewReader(testdata.Step))
	)
	assert.Nil(t, err)
	assert.False(t, actual.IsAction())
	assert.False(t, actual.IsStdoutParsable())

	// assert.Equal(t, expected, actual)
}

func TestUsesNewStepFromString(t *testing.T) {
	var (
		// expected = &sequence.Step{}
		actual, err = sequence.NewStepFromString(string(testdata.Uses))
	)
	assert.Nil(t, err)
	assert.True(t, actual.IsAction())
	assert.True(t, actual.IsStdoutParsable())

	// assert.Equal(t, expected, actual)
}


func TestUsesNewStepFromBytes(t *testing.T) {
	var (
		// expected = &sequence.Step{}
		actual, err = sequence.NewStepFromBytes(testdata.Uses)
	)
	assert.Nil(t, err)
	assert.True(t, actual.IsAction())
	assert.True(t, actual.IsStdoutParsable())

	// assert.Equal(t, expected, actual)
}

func TestUsesNewStepFromReader(t *testing.T) {
	var (
		// expected = &sequence.Step{}
		actual, err = sequence.NewStepFromReader(bytes.NewReader(testdata.Uses))
	)
	assert.Nil(t, err)
	assert.True(t, actual.IsAction())
	assert.True(t, actual.IsStdoutParsable())

	// assert.Equal(t, expected, actual)
}