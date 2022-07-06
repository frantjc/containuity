package sequence_test

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/runtime/docker"
	"github.com/frantjc/sequence/srv"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	ctx = context.TODO()
)

type RuntimeTest func(*testing.T, runtime.Runtime)

const (
	alpineRef = "docker.io/library/alpine"
	alpineImg = sequence.Image(alpineRef)

	golang118Ref = "docker.io/library/golang:1.18"
)

func NewTestRuntimes(t *testing.T) []runtime.Runtime {
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
		assert.Error(t, http.Serve(l, h2c.NewHandler(sqncRuntimeHandler, &http2.Server{})))
	}()
	t.Cleanup(func() {
		assert.Nil(t, l.Close())
	})

	time.Sleep(time.Second * 5)

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

	return []runtime.Runtime{
		dockerRuntime,
		sqncRuntime,
	}
}

func PruneTest(t *testing.T, rt runtime.Runtime) {
	assert.Nil(t, rt.PruneContainers(ctx))
	assert.Nil(t, rt.PruneVolumes(ctx))
	assert.Nil(t, rt.PruneImages(ctx))
}
