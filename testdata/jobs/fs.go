package jobs

import _ "embed"

var (
	//go:embed checkout_test_job.yml
	CheckoutTest []byte

	//go:embed env_job.yml
	Env []byte

	//go:embed save_state_job.yml
	SaveState []byte

	//go:embed set_output_job.yml
	SetOutput []byte
)
