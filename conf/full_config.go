package conf

import "os"

func NewFull(base *Config, configFilePath, repository string) (*Config, error) {
	configOpts := []ConfigOpt{
		WithConfig(base),
	}

	if configFilePath != "" {
		configOpts = append(configOpts, WithConfigFilePath(repository, configFilePath))
	}

	configOpts = append(configOpts, WithConfigFromEnv)

	if _, err := os.Stat(DefaultUserConfigFilePath); err == nil {
		configOpts = append(configOpts, WithConfigFilePath(repository, DefaultUserConfigFilePath))
	}

	configOpts = append(configOpts, WithDefaultUserConfig)

	if _, err := os.Stat(DefaultSystemConfigFilePath); err == nil {
		configOpts = append(configOpts, WithConfigFilePath(repository, DefaultSystemConfigFilePath))
	}

	configOpts = append(configOpts, WithDefaultSystemConfig)

	return New(configOpts...)
}
