package workflows

import _ "embed"

var (
	//go:embed checkout_test_build_workflow.yml
	CheckoutTestBuildWorkflow []byte

	//go:embed demo_workflow.yml
	DemoWorkflow []byte
)
