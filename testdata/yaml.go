package testdata

import _ "embed"

var (
	//go:embed action.yml
	Action []byte

	//go:embed checkout_step.yml
	CheckoutStep []byte

	//go:embed default_image_step.yml
	DefaultImageStep []byte

	//go:embed env_step.yml
	EnvStep []byte

	//go:embed checkout_test_build_workflow.yml
	CheckoutTestBuildWorkflow []byte

	//go:embed checkout_test_job.yml
	CheckoutTestJob []byte
)
