package docker

import (
	"context"
	"os"

	"github.com/docker/cli/cli/streams"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/frantjc/sequence/runtime"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/rs/zerolog/log"
)

func (r *dockerRuntime) Pull(ctx context.Context, ref string, opts ...runtime.PullOpt) (runtime.Image, error) {
	log.Debug().Msgf("pulling %s", ref)
	pref, err := name.ParseReference(ref)
	if err != nil {
		return nil, err
	}

	o, err := r.client.ImagePull(ctx, pref.Name(), types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}

	return nil, jsonmessage.DisplayJSONMessagesToStream(o, streams.NewOut(os.Stdout), nil)
}
