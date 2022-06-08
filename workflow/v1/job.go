package workflowv1

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func NewJobFromFile(name string) (*Job, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return NewJobFromReader(f)
}

func NewJobFromReader(r io.Reader) (*Job, error) {
	j := &Job{}
	d := yaml.NewDecoder(r)
	return j, d.Decode(j)
}

func (j *Job) GetStep(id string) (*Step, error) {
	for _, step := range j.Steps {
		if step.GetID() == id {
			return step, nil
		}
	}
	return nil, fmt.Errorf("job has no step with name or id %s", id)
}
