package orchestrator

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/env"
	"github.com/frantjc/sequence/meta"
	"github.com/frantjc/sequence/plan"
	"github.com/frantjc/sequence/runtime"
)

// TODO should we be using runtime.SpecOpts here?
func RunStep(ctx context.Context, r runtime.Runtime, s *sequence.Step, opts ...runtime.SpecOpt) error {
	spec, err := plan.PlanStep(ctx, s)
	if err != nil {
		return err
	}

	// TODO always pull; this conditional is for dev purposes ONLY
	if spec.Image != meta.Image() {
		_, err = r.Pull(ctx, spec.Image)
		if err != nil {
			return err
		}
	}

	for _, mount := range spec.Mounts {
		if mount.Type == runtime.MountTypeBind {
			err = os.MkdirAll(mount.Source, 0777)
			if err != nil {
				return err
			}
		}
	}

	copts := append([]runtime.SpecOpt{runtime.WithSpec(spec)}, opts...)
	container, err := r.Create(ctx, copts...)
	if err != nil {
		return err
	}

	var (
		eopts = []runtime.ExecOpt{runtime.WithStdio}
		buf   = new(bytes.Buffer)
	)
	if s.IsStdoutResponse() {
		eopts = []runtime.ExecOpt{runtime.WithStreams(os.Stdin, buf, os.Stderr)}
	}
	err = container.Exec(ctx, eopts...)
	if err != nil {
		return err
	}

	if s.IsStdoutResponse() {
		resp := &sequence.StepResponse{}
		if err = json.NewDecoder(buf).Decode(resp); err != nil {
			return err
		}

		if resp.Step != nil {
			return RunStep(ctx, r, resp.Step.Merge(s), append(opts, runtime.WithMounts(spec.Mounts...), runtime.WithEnv(env.ArrToMap(spec.Env)))...)
		}
	}

	return nil
}
