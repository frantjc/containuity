package workflows

import _ "embed"

var (
	//go:embed checkout_test_build_workflow.yml
	CheckoutTestBuild []byte

	//go:embed demo_workflow.yml
	Demo []byte

	//go:embed env_workflow.yml
	Env []byte

	//go:embed slow_workflow.yml
	Slow []byte

	//go:embed stop_commands_workflow.yml
	StopCommands []byte
)
