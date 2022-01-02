package sequence

type StepResponse struct {
	Metadata map[string]string      `json:",omitempty"`
	Version  map[string]interface{} `json:",omitempty"`
	Step     *Step                  `json:",omitempty"`
}
