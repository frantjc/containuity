package conf

import (
	"os"

	"github.com/frantjc/sequence/internal/conf/flags"
)

func NewFromFlagsWithRepository(repository string, opts ...ConfigOpt) (*Config, error) {
	configOpts := []ConfigOpt{
		WithConfig(&Config{
			Verbose:  flags.FlagVerbose,
			Port:     flags.FlagPort,
			Socket:   flags.FlagSocket,
			RootDir:  flags.FlagRootDir,
			StateDir: flags.FlagStateDir,
		}),
	}

	if flags.FlagConfigFilePath != "" {
		configOpts = append(configOpts, WithConfigFilePath(repository, flags.FlagConfigFilePath))
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
	configOpts = append(configOpts, opts...)

	return New(configOpts...)
}

func NewFromFlags(opts ...ConfigOpt) (*Config, error) {
	return NewFromFlagsWithRepository(flags.FlagWorkDir, opts...)
}
