package containerd

import (
	"context"
	"os"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cmd/ctr/commands/content"
	"github.com/frantjc/sequence/defaults"
	"github.com/frantjc/sequence/runtime"
	"github.com/google/go-containerregistry/pkg/name"
)

func (r *containerdRuntime) Pull(ctx context.Context, ref string, opts ...runtime.PullOpt) (runtime.Image, error) {
	_, err := runtime.NewPull(opts...)
	if err != nil {
		return nil, err
	}

	pref, err := name.ParseReference(ref)
	if err != nil {
		return nil, err
	}

	done := make(chan struct{}, 1)
	go func() {
		content.ShowProgress(ctx, content.NewJobs(pref.Name()), r.client.ContentStore(), os.Stdout)
		close(done)
	}()

	img, err := r.client.ImageService().Get(ctx, pref.Name())
	if err != nil {
		img, err = r.client.Fetch(ctx, pref.Name(), containerd.WithPullLabels(defaults.Labels))
		if err != nil {
			return nil, err
		}
		<-done
	}

	return nil, containerd.NewImage(r.client, img).Unpack(ctx, containerd.DefaultSnapshotter)
}
