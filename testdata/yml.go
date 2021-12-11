package testdata

import _ "embed"

//go:embed step.yml
var Step []byte

//go:embed uses.yml
var Uses []byte

//go:embed job.yml
var Job []byte

//go:embed workflow.yml
var Workflow []byte

//go:embed action.yml
var Action []byte
