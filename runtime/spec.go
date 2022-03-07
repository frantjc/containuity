package runtime

import (
	"github.com/frantjc/sequence/internal/env"
	"github.com/opencontainers/runtime-spec/specs-go"
)

type Spec struct {
	Image      string
	Entrypoint []string
	Cmd        []string
	Cwd        string
	Env        []string
	Mounts     []specs.Mount
	Privileged bool
}

type SpecOpt func(*Spec) error

func WithImage(image string) SpecOpt {
	return func(s *Spec) error {
		s.Image = image
		return nil
	}
}

func WithEntrypoint(entrypoint ...string) SpecOpt {
	return func(s *Spec) error {
		s.Entrypoint = entrypoint
		return nil
	}
}

func WithCmd(cmd ...string) SpecOpt {
	return func(s *Spec) error {
		s.Cmd = cmd
		return nil
	}
}

func WithCwd(cwd string) SpecOpt {
	return func(s *Spec) error {
		s.Cwd = cwd
		return nil
	}
}

func WithEnv(m map[string]string) SpecOpt {
	return func(s *Spec) error {
		s.Env = append(s.Env, env.MapToArr(m)...)
		return nil
	}
}

func WithEnvVar(k, v string) SpecOpt {
	return func(s *Spec) error {
		s.Env = append(s.Env, env.ToArr(k, v)...)
		return nil
	}
}

func WithPrivileged(s *Spec) error {
	s.Privileged = true
	return nil
}

func WithMounts(mounts ...specs.Mount) SpecOpt {
	return func(s *Spec) error {
		s.Mounts = append(s.Mounts, mounts...)
		return nil
	}
}
