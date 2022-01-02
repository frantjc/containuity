package container

import "fmt"

type Spec struct {
	Image      string
	Entrypoint []string
	Cmd        []string
	Cwd        string
	Env        []string
	Mounts     []Mount
	Privileged bool
}

type SpecOpt func(*Spec)

func WithImage(image string) SpecOpt {
	return func(s *Spec) {
		s.Image = image
	}
}

func WithEntrypoint(entrypoint ...string) SpecOpt {
	return func(s *Spec) {
		s.Entrypoint = entrypoint
	}
}

func WithCmd(cmd ...string) SpecOpt {
	return func(s *Spec) {
		s.Cmd = cmd
	}
}

func WithCwd(cwd string) SpecOpt {
	return func(s *Spec) {
		s.Cwd = cwd
	}
}

func WithEnv(env map[string]string) SpecOpt {
	return func(s *Spec) {
		for k, v := range env {
			s.Env = append(s.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}
}

func WithEnvVar(k, v string) SpecOpt {
	return func(s *Spec) {
		s.Env = append(s.Env, fmt.Sprintf("%s=%s", k, v))
	}
}

func WithPrivileged(s *Spec) {
	s.Privileged = true
}

func WithMounts(mounts ...Mount) SpecOpt {
	return func(s *Spec) {
		s.Mounts = append(s.Mounts, mounts...)
	}
}

func New(opts ...SpecOpt) *Spec {
	s := &Spec{}
	for _, opt := range opts {
		opt(s)
	}

	return s
}
