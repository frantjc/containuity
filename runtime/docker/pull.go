package docker

import (
	"context"

	"github.com/docker/cli/cli/streams"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/frantjc/sequence/runtime"
	"github.com/google/go-containerregistry/pkg/name"
)

func (r *dockerRuntime) Pull(ctx context.Context, ref string, opts ...runtime.PullOpt) (runtime.Image, error) {
	p, err := runtime.NewPull(opts...)
	if err != nil {
		return nil, err
	}

	pref, err := name.ParseReference(ref)
	if err != nil {
		return nil, err
	}

	o, err := r.client.ImagePull(ctx, pref.Name(), types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}

	return nil, jsonmessage.DisplayJSONMessagesToStream(o, streams.NewOut(p.Stdout), nil)
}
