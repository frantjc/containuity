package workflow

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/runtime"
	"github.com/opencontainers/runtime-spec/specs-go"
)

var (
	readOnly = []string{runtime.MountOptReadOnly}
)

const (
	containerWorkdir = "/sqnc"
	crtsDir          = "/etc/ssl"
)

func expandStep(s *Step, gctx *actions.GlobalContext) (*Step, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	expandedStep := &Step{}
	err = json.Unmarshal(
		actions.ExpandBytes(b, func(s string) string {
			return gctx.Get(s)
		}),
		expandedStep,
	)
	if err != nil {
		return nil, err
	}

	if err = actions.WithInputs(expandedStep.With)(gctx); err != nil {
		return nil, err
	}

	return expandedStep, nil
}

func runSpec(ctx context.Context, r runtime.Runtime, s *runtime.Spec, ro *runOpts, opts []runtime.ExecOpt) error {
	ro.logout.Infof("[%sSQNC%s] pulling image '%s'", log.ColorInfo, log.ColorNone, s.Image)
	image, err := r.PullImage(ctx, s.Image)
	if err != nil {
		return err
	}
	ro.logout.Debugf("[%sSQNC:DBG%s] finished pulling image '%s'", log.ColorDebug, log.ColorNone, image.Ref())

	for _, opt := range ro.specOpts {
		if err = opt(s); err != nil {
			return err
		}
	}

	container, err := r.CreateContainer(ctx, s)
	if err != nil {
		return err
	}

	return container.Exec(ctx, opts...)
}

func createFile(name string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(name), 0777); err != nil {
		return nil, err
	}

	if fs, err := os.Stat(name); err == nil && !fs.IsDir() {
		return os.Open(name)
	}

	return os.Create(name)
}

func getID(ro *runOpts) string {
	return base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprint(ro.repository, ro.workflow.Name, ro.jobName)),
	)
}

func getHostWorkdir(id string, ro *runOpts) string {
	return filepath.Join(ro.workdir, id)
}

func getHostGitHubPathFilepath(hostWorkdir string) string {
	return filepath.Join(hostWorkdir, "github", "path")
}

func getHostGitHubEnvFilepath(hostWorkdir string) string {
	return filepath.Join(hostWorkdir, "github", "env")
}

func getHostWorkspace(hostWorkdir string) string {
	return filepath.Join(hostWorkdir, "workspace")
}

func getHostRunnerTemp(hostWorkdir string) string {
	return filepath.Join(hostWorkdir, "runner", "temp")
}

func getHostRunnerToolCache(hostWorkdir string) string {
	return filepath.Join(hostWorkdir, "runner", "tool_cache")
}

func getHostActionPath(action actions.Reference, ro *runOpts) string {
	return filepath.Join(ro.workdir, "actions", action.Owner(), action.Repository(), action.Path(), action.Version())
}

func getDefaultEnv() []string {
	return []string{
		"SEQUENCE=true",
		"SQNC=true",
		"DEBIAN_FRONTEND=noninteractive",
		"ACCEPT_EULA=Y",
	}
}

func getDefaultMounts(ctx *actions.GlobalContext, stepWorkdir string) []specs.Mount {
	return []specs.Mount{
		{
			Source:      crtsDir,
			Destination: crtsDir,
			Type:        runtime.MountTypeBind,
			Options:     readonly,
		},
		{
			Source:      getHostWorkspace(stepWorkdir),
			Destination: ctx.GitHubContext.Workspace,
			Type:        runtime.MountTypeBind,
		},
		{
			Source:      getHostRunnerTemp(stepWorkdir),
			Destination: ctx.RunnerContext.Temp,
			Type:        runtime.MountTypeBind,
		},
		{
			Source:      getHostRunnerToolCache(stepWorkdir),
			Destination: ctx.RunnerContext.ToolCache,
			Type:        runtime.MountTypeBind,
		},
	}
}

func getDefaultSpec(gctx *actions.GlobalContext, stepWorkdir string, privileged bool, ro *runOpts) *runtime.Spec {
	return &runtime.Spec{
		Image:      ro.runnerImage,
		Cwd:        gctx.GitHubContext.Workspace,
		Privileged: privileged,
		Env: append(
			gctx.Arr(),
			getDefaultEnv()...,
		),
		Mounts: getDefaultMounts(gctx, stepWorkdir),
	}
}
