package sequence

import "github.com/frantjc/sequence/github/actions"

type StepResponse struct {
	Metadata map[string]string      `json:",omitempty"`
	Version  map[string]interface{} `json:",omitempty"`
	Action   *actions.Action        `json:",omitempty"`
}
