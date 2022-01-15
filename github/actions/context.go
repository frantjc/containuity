package actions

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Contexts struct {
	ctx           context.Context
	GitHubContext *GitHubContext
	EnvContext    map[string]string
	JobContext    *JobContext
	StepsContext  map[string]StepContext
	RunnerContext *RunnerContext
	InputsContext map[string]string
}

var (
	_ context.Context = &Contexts{}
)

func (c *Contexts) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c *Contexts) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Contexts) Err() error {
	return c.ctx.Err()
}

func (c *Contexts) Value(i interface{}) interface{} {
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
						return c.GitHubContext.Action
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
// https://docs.github.com/en/actions/learn-github-actions/contexts#github-context
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
	Container *Container
	Services  map[string]Service
	Status    string
}

type Container struct {
	ID      string
	Network string
}

type Service struct {
	ID      string
	Network string
	// not sure if this is correct
	Ports map[string]string
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
