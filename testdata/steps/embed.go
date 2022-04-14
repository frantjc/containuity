package steps

import _ "embed"

var (
	//go:embed checkout_step.yml
	CheckoutStep []byte

	//go:embed default_image_step.yml
	DefaultImageStep []byte

	//go:embed env_step.yml
	EnvStep []byte

	//go:embed expand_step.yml
	ExpandStep []byte

	//go:embed fail_step.yml
	FailStep []byte

	//go:embed stop_commands_step.yml
	StopCommandsStep []byte
)
