package actions

type Step struct {
	Shell            string            `json:"shell,omitempty" yaml:"shell,omitempty"`
	If               string            `json:"if,omitempty" yaml:"if,omitempty"`
	Name             string            `json:"name,omitempty" yaml:"name,omitempty"`
	ID               string            `json:"id,omitempty" yaml:"id,omitempty"`
	Env              map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	WorkingDirectory string            `json:"working-directory,omitempty" yaml:"working-directory,omitempty"`
	Uses             string            `json:"uses,omitempty" yaml:"uses,omitempty"`
	With             map[string]string `json:"with,omitempty" yaml:"with,omitempty"`
	Run              string            `json:"run,omitempty" yaml:"run,omitempty"`
}
