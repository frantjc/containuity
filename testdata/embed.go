package testdata

import (
	"github.com/frantjc/sequence/testdata/actions/composite"
	"github.com/frantjc/sequence/testdata/actions/docker"
	"github.com/frantjc/sequence/testdata/steps"
	"github.com/frantjc/sequence/testdata/jobs"
	"github.com/frantjc/sequence/testdata/workflows"
)

var (
	CompositeAction = composite.Action
	DockerAction = docker.Action
)

var (
	CheckoutStep = steps.CheckoutStep
	DefaultImageStep = steps.DefaultImageStep
	EnvStep = steps.EnvStep
	ExpandStep = steps.ExpandStep
	StopCommandsStep = steps.StopCommandsStep
)

var (
	CheckoutTestJob = jobs.CheckoutTestJob
	EnvJob = jobs.EnvJob
	SaveStateJob = jobs.SaveStateJob
	SetOutputJob = jobs.SetOutputJob
)

var (
	CheckoutTestBuildWorkflow = workflows.CheckoutTestBuildWorkflow
)
