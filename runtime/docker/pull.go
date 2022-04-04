package docker

import (
	"context"

	"github.com/docker/cli/cli/streams"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/frantjc/sequence/internal/sio"
	"github.com/frantjc/sequence/runtime"
)

type noOpWriter struct{}

func (w *noOpWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (r *dockerRuntime) PullImage(ctx context.Context, ref string) (runtime.Image, error) {
	pref, err := reference.Parse(ref)
	if err != nil {
		return nil, err
	}

	o, err := r.client.ImagePull(ctx, pref.String(), types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	defer o.Close()

	return &dockerImage{
		ref: pref.String(),
	}, jsonmessage.DisplayJSONMessagesToStream(o, streams.NewOut(sio.NewNoOpWriter()), nil)
}
