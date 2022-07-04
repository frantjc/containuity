package sequence_test

import (
	"context"
	"testing"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/runtime"
	"github.com/stretchr/testify/assert"
)

var (
	ctx = context.TODO()
)

type RuntimeTest func(*testing.T, runtime.Runtime)

const (
	alpineRef = "docker.io/library/alpine"
	alpineImg = sequence.Image(alpineRef)

	golang118Ref = "docker.io/library/golang:1.18"
	// golang118Image = sequence.Image(golang118Ref)
)

func PruneTest(t *testing.T, rt runtime.Runtime) {
	assert.Nil(t, rt.PruneContainers(ctx))
	assert.Nil(t, rt.PruneVolumes(ctx))
	assert.Nil(t, rt.PruneImages(ctx))
}
