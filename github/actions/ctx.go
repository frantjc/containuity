package actions

import (
	"context"
	"fmt"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/frantjc/sequence/github"
	"github.com/go-git/go-git/v5"
)

type ActionsContext struct {
	ctx           context.Context
	GitHubContext *GitHubContext
	EnvContext    map[string]string
	JobContext    *JobContext
	StepsContext  map[string]StepContext
	RunnerContext *RunnerContext
	InputsContext map[string]string
}

var (
	_ context.Context = &ActionsContext{}
)

func (c *ActionsContext) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c *ActionsContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *ActionsContext) Err() error {
	return c.ctx.Err()
}

func (c *ActionsContext) Value(i interface{}) interface{} {
	if s, ok := i.(string); ok {
		ss := strings.Split(s, ".")
		if len(ss) > 0 {
			switch ss[0] {
			case "github":
				if len(ss) > 1 {
					switch ss[1] {
					case "action":
						return c.GitHubContext.Action
					case "action_path":
						return c.GitHubContext.ActionPath
					case "actor":
						return c.GitHubContext.Actor
					case "base_ref":
						return c.GitHubContext.BaseRef
					case "event":
						return c.GitHubContext.Event
					case "event_name":
						return c.GitHubContext.EventName
					case "event_path":
						return c.GitHubContext.EventPath
					case "head_ref":
						return c.GitHubContext.HeadRef
					case "job":
						return c.GitHubContext.Job
					case "ref":
						return c.GitHubContext.Ref
					case "ref_name":
						return c.GitHubContext.RefName
					case "ref_protected":
						return fmt.Sprint(c.GitHubContext.RefProtected)
					case "ref_type":
						return c.GitHubContext.RefType.String()
					case "repository":
						return c.GitHubContext.Repository
					case "repository_owner":
						return c.GitHubContext.RepositoryOwner
					case "run_id":
						return c.GitHubContext.RunID
					case "run_number":
						return c.GitHubContext.RunNumber
					case "run_attempt":
						return c.GitHubContext.RunAttempt
					case "server_url":
						return c.GitHubContext.ServerURL.String()
					case "sha":
						return c.GitHubContext.Sha
					case "token":
						return c.GitHubContext.Token
					case "workflow":
						return c.GitHubContext.Workflow
					case "workspace":
						return c.GitHubContext.Workspace
					}
				}
			case "env":
				if len(ss) > 1 {
					if v, ok := c.EnvContext[ss[1]]; ok {
						return v
					}
				}
			case "job":
				if len(ss) > 1 {
					switch ss[1] {
					case "container":
						if len(ss) > 2 {
							switch ss[2] {
							case "id":
								return c.JobContext.Container.ID
							case "network":
								return c.JobContext.Container.Network
							}
						}
					case "services":
						if len(ss) > 2 {
							if v, ok := c.JobContext.Services[ss[2]]; ok {
								if len(ss) > 3 {
									switch ss[3] {
									case "id":
										return v.ID
									case "network":
										return v.Network
									case "ports":
										if len(ss) > 4 {
											if vv, ok := v.Ports[ss[4]]; ok {
												return vv
											}
										}
									}
								}
							}
						}
					case "status":
						return c.JobContext.Status
					}
				}
			case "steps":
				if len(ss) > 1 {
					if v, ok := c.StepsContext[ss[1]]; ok {
						if len(ss) > 2 {
							switch ss[2] {
							case "outputs":
								if len(ss) > 3 {
									if vv, ok := v.Outputs[ss[3]]; ok {
										return vv
									}
								}
							case "outcome":
								return v.Outcome
							case "conclusion":
								return v.Conclusion
							}
						}
					}
				}
			case "runner":
				if len(ss) > 1 {
					switch ss[1] {
					case "name":
						return c.RunnerContext.Name
					case "os":
						return c.RunnerContext.OS.String()
					case "arch":
						return c.RunnerContext.Arch.String()
					case "temp":
						return c.RunnerContext.Temp
					case "tool_cache":
						return c.RunnerContext.ToolCache
					}
				}
			case "inputs":
				if len(ss) > 1 {
					if v, ok := c.InputsContext[ss[1]]; ok {
						return v
					}
				}
			}
		}
	}

	return c.ctx.Value(i)
}

// Context represents the GitHub Context
// https://docs.github.com/en/actions/learn-github-actions/ActionsContext#github-context
type GitHubContext struct {
	Action          string
	ActionPath      string
	Actor           string
	BaseRef         string
	Event           string
	EventName       string
	EventPath       string
	HeadRef         string
	Job             string
	Ref             string
	RefName         string
	RefProtected    bool
	RefType         RefType
	Repository      string
	RepositoryOwner string
	RunID           string
	RunNumber       string
	RunAttempt      string
	ServerURL       *url.URL
	Sha             string
	Token           string
	Workflow        string
	Workspace       string
}

type JobContext struct {
	Container *struct {
		ID      string
		Network string
	}
	Services map[string]struct {
		ID      string
		Network string
		// not sure if this is the correct representation
		Ports map[string]string
	}
	Status string
}

type StepContext struct {
	Outputs    map[string]string
	Conclusion string
	Outcome    string
}

type RunnerContext struct {
	Name      string
	OS        OS
	Arch      Arch
	Temp      string
	ToolCache string
}

func defaultCtx() *ActionsContext {
	return &ActionsContext{
		GitHubContext: &GitHubContext{
			ServerURL: github.DefaultURL,
		},
		EnvContext:   map[string]string{},
		JobContext:   &JobContext{},
		StepsContext: map[string]StepContext{},
		RunnerContext: &RunnerContext{
			OS:   OSFrom(runtime.GOOS),
			Arch: ArchFrom(runtime.GOARCH),
		},
		InputsContext: map[string]string{},
	}
}

func NewContextFromPath(path string, opts ...CtxOpt) (*ActionsContext, error) {
	copts := defaultCtxOpts()
	for _, opt := range opts {
		err := opt(copts)
		if err != nil {
			return nil, err
		}
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	return newCtxFromRepository(repo, copts)
}

func newCtxFromRepository(r *git.Repository, opts *ctxOpts) (*ActionsContext, error) {
	c := defaultCtx()

	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	c.GitHubContext.Sha = ref.Hash().String()
	c.GitHubContext.RefName = ref.String()
	c.GitHubContext.Ref = ref.String()

	if ref.Name().IsBranch() {
		opts.branch = ref.String()
		c.GitHubContext.RefType = RefTypeBranch
	} else {
		c.GitHubContext.RefType = RefTypeTag
	}

	if conf, err := r.Config(); err == nil {
		c.GitHubContext.Actor = conf.Author.Name
		for _, remote := range conf.Remotes {
			for _, rurl := range remote.URLs {
				prurl, err := url.Parse(rurl)
				if err == nil {
					c.GitHubContext.Repository = strings.TrimSuffix(
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

		c.GitHubContext.RefName = branch.Name
		c.GitHubContext.Ref = branch.Name
		c.GitHubContext.RefType = RefTypeBranch
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

	return c, nil
}
