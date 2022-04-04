package actions

type Step struct {
	Shell            string            `json:"shell,omitemtpy" yaml:"shell,omitempty"`
	If               interface{}       `json:"if,omitemtpy" yaml:"if,omitempty"`
	Name             string            `json:"name,omitemtpy" yaml:"name,omitempty"`
	ID               string            `json:"id,omitemtpy" yaml:"id,omitempty"`
	Env              map[string]string `json:"env,omitemtpy" yaml:"env,omitempty"`
	WorkingDirectory string            `json:"working-directory,omitemtpy" yaml:"working-directory,omitempty"`
	Uses             string            `json:"uses,omitemtpy" yaml:"uses,omitempty"`
	With             map[string]string `json:"with,omitemtpy" yaml:"with,omitempty"`
	Run              string            `json:"run,omitempty" yaml:"run,omitempty"`
}
