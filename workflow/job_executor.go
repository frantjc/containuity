package workflow

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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
			actions.WithEnv(j.Env),
			actions.WithJobName(j.Name),
			func(gc *actions.GlobalContext) error {
				if gc.GitHubContext == nil {
					gc.GitHubContext = &actions.GitHubContext{}
				}
				if gc.RunnerContext == nil {
					gc.RunnerContext = &actions.RunnerContext{}
				}
				gc.GitHubContext.ActionPath = containerActionPath
				gc.GitHubContext.Workspace = containerWorkspace
				gc.RunnerContext.Temp = containerRunnerTemp
				gc.RunnerContext.ToolCache = containerRunnerToolCache
				return nil
			},
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

	// the path or url to the GitHub repository the actions should
	// execute with the context of
	// also used to generate a unique id so that steps of the
	// same job share state
	repository string

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
	jobEnv map[string]string
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
	e.globalContext.EnvContext = e.jobEnv
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

func (e *jobExecutor) id() string {
	rxp := regexp.MustCompile("[^a-zA-Z0-9_.-]")
	ids := []string{e.repository}
	if e.globalContext.GitHubContext.Job != "" {
		ids = append(ids, e.globalContext.GitHubContext.Job)
	}
	if e.globalContext.GitHubContext.Workflow != "" {
		ids = append(ids, e.globalContext.GitHubContext.Workflow)
	}
	return fmt.Sprintf(
		"sqnc-%s",
		strings.TrimPrefix(
			rxp.ReplaceAllLiteralString(
				strings.Join(ids, "-"),
				"-",
			),
			"-",
		),
	)
}

func (e *jobExecutor) github() string {
	return strings.Join([]string{e.id(), "github"}, "-")
}

func (e *jobExecutor) githubPath() string {
	return strings.Join([]string{e.github(), "path"}, "-")
}

func (e *jobExecutor) githubEnv() string {
	return strings.Join([]string{e.github(), "env"}, "-")
}

func (e *jobExecutor) workspace() string {
	return strings.Join([]string{e.id(), "workspace"}, "-")
}

func (e *jobExecutor) runnerTemp() string {
	return strings.Join([]string{e.id(), "runner", "temp"}, "-")
}

func (e *jobExecutor) runnerToolCache() string {
	return strings.Join([]string{e.id(), "runner", "toolcache"}, "-")
}

func (e *jobExecutor) actionPath(action actions.Reference) string {
	return strings.Join([]string{"actions", action.Owner(), action.Repository(), action.Path(), action.Version()}, "-")
}

var (
	readOnly = []string{runtime.MountOptReadOnly}
)

const (
	crtsDir     = "/etc/ssl"
	hostsFile   = "/etc/hosts"
	resolveConf = "/etc/resolv.conf"
)

var (
	containerRoot            = "/sqnc"
	containerActionPath      = filepath.Join(containerRoot, "action")
	containerWorkspace       = filepath.Join(containerRoot, "workspace")
	containerRunnerTemp      = filepath.Join(containerRoot, "runner", "temp")
	containerRunnerToolCache = filepath.Join(containerRoot, "runner", "toolcache")
	containerGitHubDir       = filepath.Join(containerRoot, "github")
	containerGitHubEnvDir    = containerGitHubDir
	containerGitHubPathDir   = containerGitHubDir
	containerGitHubEnv       = filepath.Join(containerGitHubEnvDir, "env")
	containerGitHubPath      = filepath.Join(containerGitHubPathDir, "path")
	containerShimDir         = filepath.Join(containerRoot)
	containerShim            = filepath.Join(containerShimDir, shimName)
)

func (e *jobExecutor) env() []string {
	return []string{
		fmt.Sprintf("%s=%s", actions.EnvVarEnv, containerGitHubEnv),
		fmt.Sprintf("%s=%s", actions.EnvVarPath, containerGitHubPath),
	}
}

func (e *jobExecutor) mounts() []specs.Mount {
	return []specs.Mount{
		{
			Source:      e.workspace(),
			Destination: e.globalContext.GitHubContext.Workspace,
			Type:        runtime.MountTypeVolume,
		},
		{
			Source:      e.runnerTemp(),
			Destination: e.globalContext.RunnerContext.Temp,
			Type:        runtime.MountTypeVolume,
		},
		{
			Source:      e.runnerToolCache(),
			Destination: e.globalContext.RunnerContext.ToolCache,
			Type:        runtime.MountTypeVolume,
		},
		{
			Source:      e.github(),
			Destination: containerGitHubDir,
			Type:        runtime.MountTypeVolume,
		},
		// {
		// 	Destination: containerShimDir,
		// 	Type:        runtime.MountTypeTmpfs,
		// },
		// make networking stuff act more predictably for users
		{
			Source:      crtsDir,
			Destination: crtsDir,
			Type:        runtime.MountTypeBind,
			Options:     readOnly,
		},
		{
			Source:      hostsFile,
			Destination: hostsFile,
			Type:        runtime.MountTypeBind,
			Options:     readOnly,
		},
		{
			Source:      resolveConf,
			Destination: resolveConf,
			Type:        runtime.MountTypeBind,
			Options:     readOnly,
		},
	}
}
