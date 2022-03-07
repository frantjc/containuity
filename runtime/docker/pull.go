package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
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
	io.Copy(new(noOpWriter), o)

	return &dockerImage{
		ref: pref.Name(),
	}, nil
}
