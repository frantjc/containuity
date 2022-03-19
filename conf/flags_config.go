package conf

import (
	"os"

	"github.com/frantjc/sequence/conf/flags"
)

func NewFromFlags() (*Config, error) {
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
		configOpts = append(configOpts, WithConfigFilePath(flags.FlagWorkDir, flags.FlagConfigFilePath))
	}

	configOpts = append(configOpts, WithConfigFromEnv)

	if _, err := os.Stat(DefaultUserConfigFilePath); err == nil {
		configOpts = append(configOpts, WithConfigFilePath(flags.FlagWorkDir, DefaultUserConfigFilePath))
	}

	configOpts = append(configOpts, WithDefaultUserConfig)

	if _, err := os.Stat(DefaultSystemConfigFilePath); err == nil {
		configOpts = append(configOpts, WithConfigFilePath(flags.FlagWorkDir, DefaultSystemConfigFilePath))
	}

	configOpts = append(configOpts, WithDefaultSystemConfig)

	return New(configOpts...)
}
