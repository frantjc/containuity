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
	CheckoutStep = steps.Checkout
	DefaultImageStep = steps.DefaultImage
	EnvStep = steps.Env
	ExpandStep = steps.Expand
	FailStep = steps.Fail
	StopCommandsStep = steps.StopCommands
)

var (
	CheckoutTestJob = jobs.CheckoutTest
	EnvJob = jobs.Env
	SaveStateJob = jobs.SaveState
	SetOutputJob = jobs.SetOutput
)

var (
	CheckoutTestBuildWorkflow = workflows.CheckoutTestBuild
	DemoWorkflow = workflows.Demo
	EnvWorkflow = workflows.Env
	SlowWorkflow = workflows.Slow
	StopCommandsWorkflow = workflows.StopCommands
)
