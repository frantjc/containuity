package workflow

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/frantjc/sequence/conf"
	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/runtime"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func NewJobExecutor(j *Job, opts ...ExecOpt) (Executor, error) {
	ex := &jobExecutor{
		stdout: os.Stdout,
		stderr: os.Stderr,

		globalContext: actions.EmptyContext(),

		runnerImage: conf.DefaultRunnerImage,

		steps: j.Steps,

		ctxOpts: []actions.CtxOpt{
			actions.WithWorkdir(containerWorkdir),
			actions.WithEnv(j.Env),
		},

		states: map[string]map[string]string{},
	}

	for _, opt := range opts {
		if err := opt(ex); err != nil {
			return nil, err
		}
	}

	return ex, nil
}

type executable interface {
	execute(context.Context, *jobExecutor) error
	id() string
}

type Executor interface {
	Start(context.Context) error
}

type jobExecutor struct {
	runtime runtime.Runtime

	workdir string

	stdout  io.Writer
	stderr  io.Writer
	verbose bool

	globalContext *actions.GlobalContext

	runnerImage string

	steps []*Step

	pre  []executable
	main []executable
	post []executable

	// ctxOpts are used by New functions to ensure that
	// the order of ExecOpts doesn't impact the final
	// actions.GlobalContext
	ctxOpts []actions.CtxOpt

	// states are used by a parent github action to keep
	// its pre, main and post steps linked via their state
	states map[string]map[string]string

	// env is used to reset globalContext.EnvContext after each step
	env map[string]string
}

var _ Executor = &jobExecutor{}

func (e *jobExecutor) Start(ctx context.Context) error {
	for _, step := range e.steps {
		if step.IsGitHubAction() {
			githubAction := &githubActionStep{
				ID:         step.ID,
				Name:       step.Name,
				Env:        step.Env,
				Uses:       step.Uses,
				With:       step.With,
				If:         step.If,
				Privileged: step.Privileged,
			}

			if err := githubAction.execute(ctx, e); err != nil {
				return err
			}
		} else {
			e.main = append(
				e.main,
				&regularStep{
					ID:    step.ID,
					Name:  step.Name,
					Env:   step.Env,
					Shell: step.Shell,
					Run:   step.Run,
					If:    step.If,

					Image:      step.Image,
					Entrypoint: step.Entrypoint,
					Cmd:        step.Cmd,
					Privileged: step.Privileged,
				},
			)
		}
	}

	for _, pre := range e.pre {
		if err := pre.execute(ctx, e); err != nil {
			return err
		}
		e.resetContext()
	}

	for _, main := range e.main {
		if err := main.execute(ctx, e); err != nil {
			return err
		}
		e.resetContext()
	}

	for _, post := range e.post {
		if err := post.execute(ctx, e); err != nil {
			return err
		}
		e.resetContext()
	}

	return nil
}

func (e *jobExecutor) resetContext() {
	e.globalContext.InputsContext = map[string]string{}
	e.globalContext.EnvContext = e.env
}

func (e *jobExecutor) expandStringMap(s map[string]string) map[string]string {
	m := make(map[string]string, len(s))
	for k, v := range s {
		m[k] = e.expandString(v)
	}
	return m
}

func (e *jobExecutor) expandStringArr(s []string) []string {
	a := make([]string, len(s))
	for i, b := range s {
		a[i] = e.expandString(b)
	}
	return a
}

func (e *jobExecutor) expandString(s string) string {
	return string(e.expandBytes([]byte(s)))
}

func (e *jobExecutor) expandBytes(p []byte) []byte {
	return actions.ExpandBytes(p, e.globalContext.Get)
}

func (e *jobExecutor) stepWorkdir() string {
	return filepath.Join(
		e.workdir,
		base64.StdEncoding.EncodeToString(
			[]byte(
				fmt.Sprint(
					e.globalContext.GitHubContext.Workflow,
					e.globalContext.GitHubContext.Job,
				),
			),
		),
	)
}

func (e *jobExecutor) githubPathFilepath() string {
	return filepath.Join(e.stepWorkdir(), "github", "path")
}

func (e *jobExecutor) githubEnvFilepath() string {
	return filepath.Join(e.stepWorkdir(), "github", "env")
}

func (e *jobExecutor) workspace() string {
	return filepath.Join(e.stepWorkdir(), "workspace")
}

func (e *jobExecutor) runnerTemp() string {
	return filepath.Join(e.stepWorkdir(), "runner", "temp")
}

func (e *jobExecutor) runnerToolCache() string {
	return filepath.Join(e.stepWorkdir(), "runner", "toolcache")
}

func (e *jobExecutor) actionPath(action actions.Reference) string {
	return filepath.Join(e.workdir, "actions", action.Owner(), action.Repository(), action.Path(), action.Version())
}

var (
	readOnly = []string{runtime.MountOptReadOnly}
)

const (
	containerWorkdir = "/sqnc"
	crtsDir          = "/etc/ssl"
	hostsFile        = "/etc/hosts"
	resolveConf      = "/etc/resolv.conf"
)

func (e *jobExecutor) dirMounts() []specs.Mount {
	return []specs.Mount{
		{
			Source:      crtsDir,
			Destination: crtsDir,
			Type:        runtime.MountTypeBind,
			Options:     readOnly,
		},
		{
			Source:      e.workspace(),
			Destination: e.globalContext.GitHubContext.Workspace,
			Type:        runtime.MountTypeBind,
		},
		{
			Source:      e.runnerTemp(),
			Destination: e.globalContext.RunnerContext.Temp,
			Type:        runtime.MountTypeBind,
		},
		{
			Source:      e.runnerToolCache(),
			Destination: e.globalContext.RunnerContext.ToolCache,
			Type:        runtime.MountTypeBind,
		},
	}
}

func createOrOpen(name string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(name), 0777); err != nil {
		return nil, err
	}

	if _, err := os.Stat(name); err == nil {
		return os.Open(name)
	}

	return os.Create(name)
}
