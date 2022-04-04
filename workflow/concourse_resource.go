package workflow

type ConcourseResource struct {
	Privileged bool `json:"privileged,omitempty" yaml:"privileged,omitempty"`

	Get    string                 `json:"get,omitempty" yaml:"get,omitempty"`
	Put    string                 `json:"put,omitempty" yaml:"put,omitempty"`
	Params map[string]interface{} `json:"params,omitempty" yaml:"params,omitempty"`
}
