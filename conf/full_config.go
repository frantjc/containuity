package conf

import "os"

func NewFull() (*Config, error) {
	configOpts := []ConfigOpt{
		WithConfig(&Config{
			Verbose: Verbose,
			Port: Port,
			Socket: Socket,
			RootDir: RootDir,
			StateDir: StateDir,
		}),
	}

	if ConfigFilePath != "" {
		configOpts = append(configOpts, WithConfigFilePath(WorkDir, ConfigFilePath))
	}

	configOpts = append(configOpts, WithConfigFromEnv)

	if _, err := os.Stat(DefaultUserConfigFilePath); err == nil {
		configOpts = append(configOpts, WithConfigFilePath(WorkDir, DefaultUserConfigFilePath))
	}

	configOpts = append(configOpts, WithDefaultUserConfig)

	if _, err := os.Stat(DefaultSystemConfigFilePath); err == nil {
		configOpts = append(configOpts, WithConfigFilePath(WorkDir, DefaultSystemConfigFilePath))
	}

	configOpts = append(configOpts, WithDefaultSystemConfig)

	return New(configOpts...)
}
