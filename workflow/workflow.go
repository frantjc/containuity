package workflow

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

func NewWorkflowFromReader(r io.Reader) (*Workflow, error) {
	w := &Workflow{}
	d := yaml.NewDecoder(r)
	return w, d.Decode(w)
}

type Workflow struct {
	Name string         `json:"string,omitempty"`
	Jobs map[string]Job `json:"jobs,omitempty"`
}

func (w *Workflow) GetJob(name string) (*Job, error) {
	if job, ok := w.Jobs[name]; ok {
		return &job, nil
	}
	return nil, fmt.Errorf("workflow has no job with name %s", name)
}
