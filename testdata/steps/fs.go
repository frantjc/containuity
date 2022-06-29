package steps

import _ "embed"

var (
	//go:embed checkout_step.yml
	Checkout []byte

	//go:embed default_image_step.yml
	DefaultImage []byte

	//go:embed env_step.yml
	Env []byte

	//go:embed expand_step.yml
	Expand []byte

	//go:embed fail_step.yml
	Fail []byte

	//go:embed slow_step.yml
	Slow []byte

	//go:embed stop_commands_step.yml
	StopCommands []byte
)
