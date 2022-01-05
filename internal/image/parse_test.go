package image_test

import (
	"fmt"
	"testing"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/image"
	"github.com/stretchr/testify/assert"
)

func TestParseRef(t *testing.T) {
	var (
		ref         = sequence.Repository
		expected    = fmt.Sprintf("index.docker.io/%s:latest", ref)
		actual, err = image.ParseRef(ref)
	)
	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
}
