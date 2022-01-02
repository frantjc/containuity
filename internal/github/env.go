package github

import (
	"fmt"
	"net/url"
	"os"
	"runtime"

	"github.com/frantjc/sequence/internal/env"
	"github.com/go-git/go-git/v5"
)

type RefType int

const (
	RefTypeTag RefType = iota
	RefTypeBranch
)

func (r RefType) String() string {
	switch r {
	case RefTypeBranch:
		return "branch"
	case RefTypeTag:
		return "tag"
	}

	return ""
}

type OS int

const (
	LinuxOS OS = iota
	WindowsOS
	DarwinOS
)

func (o OS) String() string {
	switch o {
	case LinuxOS:
		return "Linux"
	case WindowsOS:
		return "Windows"
	case DarwinOS:
		return "macOS"
	}

	return ""
}

func OSFrom(s string) OS {
	switch s {
	case "darwin":
		return DarwinOS
	case "linux":
		return LinuxOS
	case "windows":
		return WindowsOS
	}

	return -1
}

type Arch int

const (
	X86Arch Arch = iota
	X64Arch
	ARMArch
	ARM64Arch
)

func (a Arch) String() string {
	switch a {
	case X86Arch:
		return "X86"
	case X64Arch:
		return "X64"
	case ARMArch:
		return "ARM"
	case ARM64Arch:
		return "ARM64"
	}

	return ""
}

func ArchFrom(s string) Arch {
	switch s {
	case "amd64":
		return X86Arch
	}

	return -1
}

type Env struct {
	CI              bool
	Workflow        string
	RunID           int
	RunNumber       int
	Job             string
	Action          string
	ActionPath      string
	Actions         bool
	Actor           string
	Repository      string
	EventName       string
	EventPath       string
	Workspace       string
	Sha             string
	Ref             string
	RefName         string
	RefProtected    bool
	RefType         RefType
	HeadRef         string
	BaseRef         string
	ServerURL       *url.URL
	APIURL          *url.URL
	GraphQLURL      *url.URL
	RunnerName      string
	RunnerOS        OS   // Linux, Windows or macOS
	RunnerArch      Arch // X86, X64, ARM or ARM64
	RunnerTemp      string
	RunnerToolCache string
}

func (e *Env) Map() map[string]string {
	return map[string]string{
		EnvVarCI:              fmt.Sprint(e.CI),
		EnvVarWorkflow:        e.Workflow,
		EnvVarRunID:           fmt.Sprint(e.RunID),
		EnvVarRunNumber:       fmt.Sprint(e.RunNumber),
		EnvVarJob:             e.Job,
		EnvVarAction:          e.Action,
		EnvVarActionPath:      e.ActionPath,
		EnvVarActions:         fmt.Sprint(e.Actions),
		EnvVarActor:           e.Actor,
		EnvVarRepository:      e.Repository,
		EnvVarEventName:       e.EventName,
		EnvVarEventPath:       e.EventPath,
		EnvVarWorkspace:       e.Workspace,
		EnvVarSha:             e.Sha,
		EnvVarRef:             e.Ref,
		EnvVarRefName:         e.RefName,
		EnvVarRefProtected:    fmt.Sprint(e.RefProtected),
		EnvVarRefType:         e.RefType.String(),
		EnvVarHeadRef:         e.HeadRef,
		EnvVarBaseRef:         e.BaseRef,
		EnvVarServerURL:       e.ServerURL.String(),
		EnvVarAPIURL:          e.APIURL.String(),
		EnvVarGraphQLURL:      e.GraphQLURL.String(),
		EnvVarRunnerName:      e.RunnerName,
		EnvVarRunnerOS:        e.RunnerOS.String(),
		EnvVarRunnerArch:      e.RunnerArch.String(),
		EnvVarRunnerTemp:      e.RunnerTemp,
		EnvVarRunnerToolCache: e.RunnerToolCache,
	}
}

func (e *Env) Arr() []string {
	return env.MapToArr(e.Map())
}

type envOpts struct {
	remote    string
	branch    string
	tmpdir    string
	workdir   string
	toolcache string
	runID     int
	runNumber int
}

type EnvOpt func(e *envOpts) error

func defaultEnvOpts() *envOpts {
	cachedir, _ := os.UserCacheDir()

	return &envOpts{
		remote:    "origin",
		branch:    "main",
		tmpdir:    os.TempDir(),
		toolcache: cachedir,
	}
}

func NewEnv(path string, opts ...EnvOpt) (*Env, error) {
	e := &Env{
		CI:              true,
		Actions:         true,
		RunnerOS:        OSFrom(runtime.GOOS),
		RunnerArch:      ArchFrom(runtime.GOARCH),
		ServerURL:       URL,
		APIURL:          APIURL,
		GraphQLURL:      GraphQLURL,
		Workspace:       "/tmp/sqnc/workspace",
		ActionPath:      "/tmp/sqnc/action",
		RunnerTemp:      "/tmp/sqnc/runner/temp",
		RunnerToolCache: "/tmp/sqnc/runner/toolcache",
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	e.Sha = ref.Hash().String()

	eopts := defaultEnvOpts()
	if ref.Name().IsBranch() {
		eopts.branch = ref.String()
		e.RefName = ref.String()
		e.Ref = ref.String()
		e.RefType = RefTypeBranch
	}

	for _, opt := range opts {
		err := opt(eopts)
		if err != nil {
			return nil, err
		}
	}

	if conf, err := repo.Config(); err == nil {
		e.Actor = conf.Author.Name
	}

	if branch, err := repo.Branch(eopts.branch); err == nil {
		if eopts.remote == "" {
			eopts.remote = branch.Remote
		}

		e.RefName = branch.Name
		e.Ref = branch.Name
		e.RefType = RefTypeBranch
	}

	if remote, err := repo.Remote(eopts.remote); err == nil {
		for _, u := range remote.Config().URLs {
			pu, err := url.Parse(u)
			if err == nil {
				e.ServerURL = pu
				break
			}
		}
	}

	return e, nil
}
