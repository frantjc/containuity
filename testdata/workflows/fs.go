package workflows

import _ "embed"

var (
	//go:embed checkout_test_build_workflow.yml
	CheckoutTestBuild []byte

	//go:embed demo_workflow.yml
	Demo []byte
)
