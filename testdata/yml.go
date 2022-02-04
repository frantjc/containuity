package testdata

import _ "embed"

var (
	//go:embed action.yml
	Action []byte

	//go:embed job.yml
	Job []byte

	//go:embed step.yml
	Step []byte

	//go:embed uses.yml
	Uses []byte

	//go:embed workflow.yml
	Workflow []byte
)
