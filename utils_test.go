package sequence_test

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/pkg/github/actions"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/runtime/docker"
	"github.com/frantjc/sequence/srv"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

type RuntimeTest func(*testing.T, runtime.Runtime)

const (
	alpineRef = "docker.io/library/alpine"
	alpineImg = sequence.Image(alpineRef)

	golang118Ref = "docker.io/library/golang:1.18"
)

func NewTestRuntimes(t *testing.T) []runtime.Runtime {
	t.Helper()

	dockerRuntime, err := docker.NewRuntime(ctx)
	assert.NotNil(t, dockerRuntime)
	assert.Nil(t, err)

	// listen on a random port
	l, err := net.Listen("tcp", "localhost:0")
	assert.NotNil(t, l)
	assert.Nil(t, err)

	sqncRuntimeHandler, err := srv.NewHandler(ctx, srv.WithRuntime(dockerRuntime))
	assert.NotNil(t, sqncRuntimeHandler)
	assert.Nil(t, err)

	// serve sqncRuntimeHandler on the random port
	go func() {
		assert.Error(t, http.Serve(l, sqncRuntimeHandler))
	}()
	t.Cleanup(func() {
		assert.Nil(t, l.Close())
	})

	// get the random address
	addr, err := url.Parse("http://" + l.Addr().String())
	assert.NotNil(t, addr)
	assert.Nil(t, err)

	// create a client connected to the random port
	client, err := sequence.NewClient(ctx, addr)
	assert.NotNil(t, client)
	assert.Nil(t, err)

	sqncRuntime := client.Runtime()
	assert.NotNil(t, sqncRuntime)

	runtimes := []runtime.Runtime{
		dockerRuntime,
		sqncRuntime,
	}

	for _, rt := range runtimes {
		t.Cleanup(func() {
			PruneTest(t, rt)
		})
	}

	return runtimes
}

func NewTestGlobalContext(t *testing.T) *actions.GlobalContext {
	t.Helper()

	var (
		githubToken = os.Getenv("GITHUB_TOKEN")
		// all tests in this suite are ran against
		// https://github.com/frantjc/sequence
		wd, err = os.Getwd()
	)
	assert.Nil(t, err)
	if !assert.NotEmpty(t, githubToken) {
		assert.FailNow(t, "GITHUB_TOKEN must be set")
	}

	gc, err := actions.NewContextFromPath(ctx, wd, actions.WithToken(githubToken))
	assert.NotNil(t, gc)
	assert.Nil(t, err)

	return gc
}

func PruneTest(t *testing.T, rt runtime.Runtime) {
	t.Helper()

	assert.Nil(t, rt.PruneContainers(ctx))

	if prune, err := strconv.ParseBool(os.Getenv("SQNC_PRUNE_VOLUMES")); err == nil && prune {
		assert.Nil(t, rt.PruneVolumes(ctx))
	}

	if prune, err := strconv.ParseBool(os.Getenv("SQNC_PRUNE_IMAGES")); err == nil && prune {
		assert.Nil(t, rt.PruneImages(ctx))
	}
}
