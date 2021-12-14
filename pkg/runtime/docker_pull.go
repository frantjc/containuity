package runtime

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/frantjc/sequence"
)

func (d *dockerRuntime) getPullOpts(ctx context.Context, s *sequence.Step, r *runOpts) (types.ImagePullOptions, error) {
	pullOpts := types.ImagePullOptions{}
	// TODO for private image repositories
	// creds, err := fromOpts?
	// if err != nil {
	// 	return err
	// }

	// if auth := creds.Base64(); auth != "" {
	// 	pullOpts.RegistryAuth = auth
	// }

	return pullOpts, nil
}
