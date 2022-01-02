package orchestrator

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/pkg/container"
	"github.com/frantjc/sequence/pkg/planner"
	"github.com/frantjc/sequence/pkg/runtime"
	"github.com/frantjc/sequence/pkg/sio"
	"github.com/rs/zerolog/log"
)

func RunStep(ctx context.Context, r runtime.Runtime, s *sequence.Step, opts ...RunOpt) error {
	ro := &runOpts{}
	for _, opt := range opts {
		err := opt(ro)
		if err != nil {
			log.Debug().Err(err).Msgf("planning step failed %s", s.ID())
			return err
		}
	}

	po := []planner.PlanOpt{}
	if ro.path != "" {
		po = append(po, planner.WithPath(ro.path))
	}

	spec, err := planner.PlanStep(ctx, s, po...)
	if err != nil {
		log.Debug().Err(err).Msgf("planning step failed %s", s.ID())
		return err
	}

	spec.Mounts = append(spec.Mounts, ro.mounts...)
	spec.Env = append(spec.Env, ro.env...)

	// TODO always pull, this is tmp for dev
	if spec.Image != sequence.Image() {
		err = r.Pull(ctx, spec.Image)
		if err != nil {
			log.Debug().Err(err).Msgf("pulling image failed %s", spec.Image)
			return err
		}
	}

	for _, m := range spec.Mounts {
		if m.Type == container.MountTypeBind {
			err = os.MkdirAll(m.Source, 0755)
			if err != nil {
				log.Debug().Err(err).Msgf("creating dir failed %s", m.Source)
				return err
			}
		}
	}

	c, err := r.Create(ctx, spec)
	if err != nil {
		log.Debug().Err(err).Msg("creating container failed")
		return err
	}

	streams := sio.New(sio.WithStdio)
	buf := new(bytes.Buffer)
	resp := &sequence.StepResponse{}
	if s.IsStdoutParsable() {
		streams = sio.New(sio.WithStreams(os.Stdin, buf, os.Stderr))
	}

	err = c.Start(ctx, streams)
	if err != nil {
		log.Debug().Err(err).Msg("starting container failed")
		return err
	}

	if s.IsStdoutParsable() {
		err = json.NewDecoder(buf).Decode(resp)
		if err != nil {
			log.Debug().Err(err).Msgf("parsing stdout failed %s", buf.String())
			return err
		}
	}

	if resp.Step != nil {
		resp.Step.IDF = s.IDF
		resp.Step.Name = s.Name
		return RunStep(ctx, r, resp.Step, append(opts, WithMounts(spec.Mounts), WithEnv(spec.Env))...)
	}

	return nil
}
