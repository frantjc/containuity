package docker

import (
	"context"

	"github.com/docker/cli/cli/streams"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/runtime"
	"github.com/google/go-containerregistry/pkg/name"
)

type noOpWriter struct{}

func (w *noOpWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (r *dockerRuntime) PullImage(ctx context.Context, ref string) (runtime.Image, error) {
	pref, err := name.ParseReference(ref)
	if err != nil {
		return nil, err
	}

	o, err := r.client.ImagePull(ctx, pref.Name(), types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	defer o.Close()

	return &dockerImage{
		ref: pref.Name(),
	}, jsonmessage.DisplayJSONMessagesToStream(o, streams.NewOut(log.Writer()), nil)
}
