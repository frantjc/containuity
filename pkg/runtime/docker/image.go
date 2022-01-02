package docker

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/frantjc/sequence/internal/image"
)

func (r *dockerRuntime) Pull(ctx context.Context, ref string) error {
	if r == nil {
		return fmt.Errorf("nil runtime")
	}

	pref, err := image.ParseRef(ref)
	if err != nil {
		return err
	}

	o, err := r.client.ImagePull(ctx, pref, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	// TODO write to stream provided by opts
	io.Copy(new(bytes.Buffer), o)

	return nil
}
