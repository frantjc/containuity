package workflow

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func NewWorkflowFromFile(name string) (*Workflow, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return NewWorkflowFromReader(f)
}

func NewWorkflowFromReader(r io.Reader) (*Workflow, error) {
	w := &Workflow{}
	d := yaml.NewDecoder(r)
	return w, d.Decode(w)
}

type Workflow struct {
	Name string          `json:"name,omitempty" yaml:"name,omitempty"`
	Jobs map[string]*Job `json:"jobs,omitempty" yaml:"jobs,omitempty"`
}

func (w *Workflow) GetJob(name string) (*Job, error) {
	if job, ok := w.Jobs[name]; ok {
		return job, nil
	}
	return nil, fmt.Errorf("workflow has no job with name %s", name)
}

func (w *Workflow) GetStep(id string) (*Step, error) {
	for _, job := range w.Jobs {
		if step, err := job.GetStep(id); err == nil {
			return step, nil
		}
	}
	return nil, fmt.Errorf("workflow has no step with id %s", id)
}
