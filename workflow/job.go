package workflow

import (
	"fmt"
	"io"
	"net/url"

	"gopkg.in/yaml.v3"
)

func NewJobFromReader(r io.Reader) (*Job, error) {
	j := &Job{}
	d := yaml.NewDecoder(r)
	return j, d.Decode(j)
}

type Job struct {
	Name        string            `json:"name,omitempty"`
	Permissions interface{}       `json:"permissions,omitempty"`
	Needs       interface{}       `json:"needs,omitempty"`
	If          bool              `json:"if,omitempty"`
	RunsOn      string            `json:"runs-on,omitempty"`
	Environment *Environment      `json:"environment,omitempty"`
	Concurrency interface{}       `json:"concurrency,omitempty"`
	Outputs     map[string]string `json:"outputs,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Container   interface{}       `json:"container,omitempty"`
	Steps       []Step            `json:"steps,omitempty"`
}

type Environment struct {
	Name string   `json:"name,omitempty"`
	URL  *url.URL `json:"url,omitempty"`
}

type Concurrency struct {
	Group            string `json:"group,omitempty"`
	CancelInProgress bool   `json:"cancel-in-progress,omitempty"`
}

type Defaults struct {
	Run *Run `json:"run,omitempty"`
}

type Run struct {
	Shell            string `json:"shell,omitempty"`
	WorkingDirectory string `json:"working-directory,omitempty"`
}

type Container struct {
	Image string `json:"image,omitempty"`
}

func (j *Job) GetStep(id string) (*Step, error) {
	for _, step := range j.Steps {
		if step.GetID() == id {
			return &step, nil
		}
	}
	return nil, fmt.Errorf("job has no step with name or id %s", id)
}
