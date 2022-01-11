package containerd

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/oci"
	"github.com/frantjc/sequence/defaults"
	"github.com/frantjc/sequence/runtime"
)

func (r *containerdRuntime) Create(ctx context.Context, opts ...runtime.SpecOpt) (runtime.Container, error) {
	spec, err := runtime.NewSpec(opts...)
	if err != nil {
		return nil, err
	}

	image, err := r.client.GetImage(ctx, spec.Image)
	if err != nil {
		return nil, err
	}

	ociopts := []oci.SpecOpts{
		oci.WithDefaultSpec(),
		oci.WithImageConfig(image),
		oci.WithProcessArgs(append(spec.Entrypoint, spec.Cmd...)...),
		oci.WithProcessCwd(spec.Cwd),
		oci.WithEnv(spec.Env),
		oci.WithMounts(spec.Mounts),
	}
	if spec.Privileged {
		ociopts = append(ociopts, oci.WithAddedCapabilities(privilegedCapabilities))
	}
	copts := []containerd.NewContainerOpts{
		containerd.WithNewSnapshot("", image),
		containerd.WithNewSpec(ociopts...),
		containerd.WithAdditionalContainerLabels(defaults.Labels),
	}
	container, err := r.client.NewContainer(ctx, "", copts...)
	if err != nil {
		return nil, err
	}

	return &containerdContainer{container}, nil
}
