package actions

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/pkg/container"
)

func ActionMounts(id string) []container.Mount {
	return []container.Mount{
		{
			Source:      filepath.Join("/tmp", id, "workdir"),
			Destination: filepath.Join("/tmp", id, "workdir"),
			Type:        container.MountTypeBind,
		},
		{
			Source:      filepath.Join("/tmp", id, "toolcache"),
			Destination: filepath.Join("/tmp", id, "toolcache"),
			Type:        container.MountTypeBind,
		},
		{
			Destination: filepath.Join("/tmp", id, "runnertemp"),
			Type:        container.MountTypeTmpfs,
		},
		{
			Source:      filepath.Join("/tmp", id, "action"),
			Destination: filepath.Join("/tmp", id, "action"),
			Type:        container.MountTypeBind,
		},
	}
}

var (
	mountVars = []string{
		"GITHUB_WORKSPACE",
		"RUNNER_TOOL_CACHE",
		"RUNNER_TEMP",
		"GITHUB_ACTIONS_PATH",
	}
	defaultEnv = []string{
		fmt.Sprintf("RUNNER_OS=%s", runtime.GOOS),
		fmt.Sprintf("RUNNER_ARCH=%s", runtime.GOARCH),
	}
)

func ActionEnv(id string) []string {
	m := ActionMounts(id)
	if len(m) < len(mountVars) {
		panic(fmt.Sprintf("%s/internal/actions.mountVars must be in sync with %s/internal/actions.ActionMounts", sequence.Module, sequence.Module))
	}

	s := defaultEnv
	for i, v := range mountVars {
		s = append(s, fmt.Sprintf("%s=%s", v, m[i]))
	}

	return s
}
