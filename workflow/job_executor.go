package workflow

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/internal/conf"
	"github.com/frantjc/sequence/internal/log"
	"github.com/frantjc/sequence/runtime"
	runtimev1 "github.com/frantjc/sequence/runtime/v1"
	workflowv1 "github.com/frantjc/sequence/workflow/v1"
)

func NewJobExecutor(j *workflowv1.Job, opts ...ExecOpt) (Executor, error) {
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

	steps []*workflowv1.Step

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
	logout := log.New(e.stdout).SetVerbose(e.verbose)
	logout.Infof("[%sSQNC%s] running job '%s'", log.ColorInfo, log.ColorNone, e.globalContext.GitHubContext.Job)
	for _, step := range e.steps {
		if step.IsGitHubAction() {
			githubAction := &githubActionStep{
				ID:         step.Id,
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
					ID:    step.Id,
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

var idRxp = regexp.MustCompile("[^a-zA-Z0-9_.-]")

func (e *jobExecutor) id() string {
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
			idRxp.ReplaceAllLiteralString(
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
	ids := []string{"sqnc", "actions", action.Owner(), action.Repository()}
	if action.Path() != "" {
		ids = append(ids, action.Path())
	}
	ids = append(ids, action.Version())
	return idRxp.ReplaceAllLiteralString(
		strings.Join(ids, "-"),
		"-",
	)
}

var (
	readOnly = []string{runtimev1.MountOptReadOnly}
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
	containerShimDir         = containerRoot
	containerShim            = filepath.Join(containerShimDir, shimName)
)

func (e *jobExecutor) env() []string {
	return []string{
		fmt.Sprintf("%s=%s", actions.EnvVarEnv, containerGitHubEnv),
		fmt.Sprintf("%s=%s", actions.EnvVarPath, containerGitHubPath),
	}
}

func (e *jobExecutor) mounts() []*runtimev1.Mount {
	return []*runtimev1.Mount{
		{
			Source:      e.workspace(),
			Destination: e.globalContext.GitHubContext.Workspace,
			Type:        runtimev1.MountTypeVolume,
		},
		{
			Source:      e.runnerTemp(),
			Destination: e.globalContext.RunnerContext.Temp,
			Type:        runtimev1.MountTypeVolume,
		},
		{
			Source:      e.runnerToolCache(),
			Destination: e.globalContext.RunnerContext.ToolCache,
			Type:        runtimev1.MountTypeVolume,
		},
		{
			Source:      e.github(),
			Destination: containerGitHubDir,
			Type:        runtimev1.MountTypeVolume,
		},
		// {
		// 	Destination: containerShimDir,
		// 	Type:        runtime.MountTypeTmpfs,
		// },
		// make networking stuff act more predictably for users
		{
			Source:      crtsDir,
			Destination: crtsDir,
			Type:        runtimev1.MountTypeBind,
			Options:     readOnly,
		},
		{
			Source:      hostsFile,
			Destination: hostsFile,
			Type:        runtimev1.MountTypeBind,
			Options:     readOnly,
		},
		{
			Source:      resolveConf,
			Destination: resolveConf,
			Type:        runtimev1.MountTypeBind,
			Options:     readOnly,
		},
	}
}
