package actions

type Workflow struct {
	Name string          `json:"name,omitempty" yaml:"name,omitempty"`
	Jobs map[string]*Job `json:"jobs,omitempty" yaml:"jobs,omitempty"`
}
