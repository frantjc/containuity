package sequence

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

type Workflow struct {
	Name string
	Jobs map[string]Job
}

func NewWorkflowFromBytes(b []byte) (*Workflow, error) {
	w := &Workflow{}
	return w, yaml.Unmarshal(b, w)
}

func NewWorkflowFromReader(r io.Reader) (*Workflow, error) {
	w := &Workflow{}
	d := yaml.NewDecoder(r)
	return w, d.Decode(w)
}

func NewWorkflowFromString(s string) (*Workflow, error) {
	return NewWorkflowFromBytes([]byte(s))
}

func (w *Workflow) GetJob(name string) (*Job, error) {
	if job, ok := w.Jobs[name]; ok {
		return &job, nil
	}
	return nil, fmt.Errorf("workflow has no job with name %s", name)
}
