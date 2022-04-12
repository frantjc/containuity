package docker

import (
	"context"
	"io"

	"github.com/docker/cli/cli/streams"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/frantjc/sequence/runtime"
)

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
	jsonmessage.DisplayJSONMessagesToStream(o, streams.NewOut(io.Discard), nil)

	_, _, err = r.client.ImageInspectWithRaw(ctx, pref.String())
	if err != nil {
		return nil, err
	}

	return &dockerImage{
		ref: pref.String(),
	}, nil
}
