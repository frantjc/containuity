package sequence

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

func NewJobFromBytes(b []byte) (*Job, error) {
	j := &Job{}
	return j, yaml.Unmarshal(b, j)
}

func NewJobFromReader(r io.Reader) (*Job, error) {
	j := &Job{}
	d := yaml.NewDecoder(r)
	return j, d.Decode(j)
}

func NewJobFromString(s string) (*Job, error) {
	return NewJobFromBytes([]byte(s))
}

type Job struct {
	StepsF []Step `yaml:"steps"`
}

var (
	_ Steppable = &Job{}
)

func (j *Job) Steps() ([]Step, error) {
	return j.StepsF, nil
}

func (j *Job) Step(id string) (*Step, error) {
	for _, step := range j.StepsF {
		if step.ID() == id {
			return &step, nil
		}
	}

	return nil, fmt.Errorf("job has no step with name or id %s", id)
}
