package containerd

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/frantjc/sequence/runtime"
	"github.com/google/go-containerregistry/pkg/name"
)

func (r *containerdRuntime) Pull(ctx context.Context, ref string, opts ...runtime.PullOpt) error {
	_, err := runtime.NewPull(opts...)
	if err != nil {
		return err
	}

	pref, err := name.ParseReference(ref)
	if err != nil {
		return err
	}

	img, err := r.client.ImageService().Get(ctx, pref.Name())
	if err != nil {
		img, err = r.client.Fetch(ctx, pref.Name())
		if err != nil {
			return err
		}
	}

	return containerd.NewImage(r.client, img).Unpack(ctx, containerd.DefaultSnapshotter)
}
