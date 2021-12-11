package runtime

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/frantjc/sequence"
	"github.com/google/go-containerregistry/pkg/name"
)

func (d *dockerRuntime) getContainerConfig(ctx context.Context, s *sequence.Step, r *runOpts) (*container.Config, error) {
	var (
		ref        string
		entrypoint []string
		cmd        []string
	)
	if s.Image != "" {
		ref = s.Image
		entrypoint = s.Entrypoint
		cmd = s.Cmd
	} else if s.Uses != "" {
		// TODO use special container image with github.com/frantjc/sequence/cmd/uses as the entrypoint
		ref = "alpine/git"
		entrypoint = []string{"git"}
		cmd = []string{"clone", fmt.Sprintf("https://github.com/%s", s.Uses)}
	}

	image, err := name.ParseReference(ref)
	if err != nil {
		return nil, fmt.Errorf("unable to parse image ref %s", ref)
	}

	return &container.Config{
		Image:      image.Name(),
		Entrypoint: entrypoint,
		Cmd:        cmd,
		Labels:     labels,
	}, nil
}

func (d *dockerRuntime) getHostConfig(ctx context.Context, s *sequence.Step, r *runOpts) (*container.HostConfig, error) {
	return &container.HostConfig{
		AutoRemove: true,
		Privileged: s.Privileged,
	}, nil
}
