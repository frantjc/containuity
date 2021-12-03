package testdata

import _ "embed"

//go:embed step.yml
var Step []byte

//go:embed job.yml
var Job []byte

//go:embed workflow.yml
var Workflow []byte
