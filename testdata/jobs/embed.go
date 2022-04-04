package jobs

import _ "embed"

var (
	//go:embed checkout_test_job.yml
	CheckoutTestJob []byte

	//go:embed env_job.yml
	EnvJob []byte

	//go:embed save_state_job.yml
	SaveStateJob []byte

	//go:embed set_output_job.yml
	SetOutputJob []byte
)
