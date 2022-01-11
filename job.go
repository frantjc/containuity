package sequence

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

func NewJobFromReader(r io.Reader) (*Job, error) {
	j := &Job{}
	d := yaml.NewDecoder(r)
	return j, d.Decode(j)
}

type Job struct {
	Steps []Step `json:",omitempty"`
}

func (j *Job) GetStep(id string) (*Step, error) {
	for _, step := range j.Steps {
		if step.GetID() == id {
			return &step, nil
		}
	}
	return nil, fmt.Errorf("job has no step with name or id %s", id)
}
