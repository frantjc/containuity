package actions

import (
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/frantjc/sequence/env"
	"github.com/frantjc/sequence/github"
	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
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

	Env  string
	Path string

	Token string
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

		EnvVarEnv:  e.Env,
		EnvVarPath: e.Path,

		EnvVarToken: e.Token,
	}
}

func (e *Env) Arr() []string {
	return env.MapToArr(e.Map())
}

func defaultEnv() *Env {
	return &Env{
		CI:         true,
		Actions:    true,
		ServerURL:  github.DefaultURL,
		APIURL:     github.DefaultAPIURL,
		GraphQLURL: github.DefaultGraphQLURL,
		RunnerOS:   OSFrom(runtime.GOOS),
		RunnerArch: ArchFrom(runtime.GOARCH),
	}
}

func NewEnvFromPath(path string, opts ...VarsOpt) (*Env, error) {
	vopts := defaultVarsOpts()
	for _, opt := range opts {
		err := opt(vopts)
		if err != nil {
			return nil, err
		}
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	return newEnvFromRepository(repo, vopts)
}

// get from cli flags, env, config file, .git or remote
func newEnvFromRepository(r *git.Repository, opts *varsOpts) (*Env, error) {
	e := defaultEnv()
	e.Workflow = opts.workflow
	e.RunID = opts.runID
	e.RunNumber = opts.runNumber
	e.Job = opts.job
	e.Action = uuid.NewString()
	e.ActionPath = filepath.Join(opts.workdir, "action")
	e.Workspace = filepath.Join(opts.workdir, "workspace")
	e.RefProtected = opts.refProtected
	e.HeadRef = opts.headRef
	e.BaseRef = opts.baseRef
	e.RunnerName = opts.runnerName
	e.RunnerTemp = filepath.Join(opts.workdir, "runner", "temp")
	e.RunnerToolCache = filepath.Join(opts.workdir, "runner", "toolcache")
	e.Env = filepath.Join(opts.workdir, "github", "env")
	e.Path = filepath.Join(opts.workdir, "github", "path")
	e.Token = opts.token

	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	e.Sha = ref.Hash().String()
	e.RefName = ref.String()
	e.Ref = ref.String()

	if ref.Name().IsBranch() {
		opts.branch = ref.String()
		e.RefType = RefTypeBranch
	} else {
		e.RefType = RefTypeTag
	}

	if conf, err := r.Config(); err == nil {
		e.Actor = conf.Author.Name
		for _, remote := range conf.Remotes {
			for _, rurl := range remote.URLs {
				prurl, err := url.Parse(rurl)
				if err == nil {
					e.Repository = strings.TrimSuffix(
						strings.TrimPrefix(prurl.Path, "/"),
						".git",
					)
					break
				}
			}
		}
	}

	if branch, err := r.Branch(opts.branch); err == nil {
		if opts.remote == "" {
			opts.remote = branch.Remote
		}

		e.RefName = branch.Name
		e.Ref = branch.Name
		e.RefType = RefTypeBranch
	}

	if remote, err := r.Remote(opts.remote); err == nil {
		for _, u := range remote.Config().URLs {
			_, err := url.Parse(u)
			if err == nil {
				// override default github urls
				break
			}
		}
	}

	return e, nil
}
