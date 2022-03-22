package conf

import (
	"io"
	"net/url"
	"os"
	"strconv"

	"github.com/pelletier/go-toml/v2"
)

type ConfigOpt func() (*Config, error)

func WithDefaultSystemConfig() (*Config, error) {
	return &Config{
		Socket:   DefaultSystemSocket,
		RootDir:  DefaultSystemRootDir,
		StateDir: DefaultSystemStateDir,
		GitHub: &GitHubConfig{
			URL: DefaultGitHubURL,
		},
		Runtime: &RuntimeConfig{
			Name:        DefaultRuntimeName,
			RunnerImage: DefaultRunnerImage,
		},
	}, nil
}

func WithDefaultUserConfig() (*Config, error) {
	return &Config{
		Socket:   DefaultUserSocket,
		RootDir:  DefaultUserRootDir,
		StateDir: DefaultUserStateDir,
		GitHub: &GitHubConfig{
			URL: DefaultGitHubURL,
		},
		Runtime: &RuntimeConfig{
			Name:        DefaultRuntimeName,
			RunnerImage: DefaultRunnerImage,
		},
	}, nil
}

func WithConfigFromEnv() (c *Config, err error) {
	c = &Config{
		Socket:   os.Getenv(EnvVarSocket),
		RootDir:  os.Getenv(EnvVarRootDir),
		StateDir: os.Getenv(EnvVarStateDir),
		GitHub: &GitHubConfig{
			Token: os.Getenv(EnvVarGitHubToken),
		},
		Runtime: &RuntimeConfig{
			Name:        os.Getenv(EnvVarRuntime),
			RunnerImage: os.Getenv(EnvVarRunnerImage),
			ActionImage: os.Getenv(EnvVarActionImage),
		},
		Secrets: map[string]string{},
	}
	if rawVerbose := os.Getenv(EnvVarVerbose); rawVerbose != "" {
		c.Verbose, err = strconv.ParseBool(rawVerbose)
		if err != nil {
			return
		}
	}

	if rawPort := os.Getenv(EnvVarPort); rawPort != "" {
		c.Port, err = strconv.Atoi(rawPort)
		if err != nil {
			return
		}
	}

	if rawGitHubURL := os.Getenv(EnvVarGitHubURL); rawGitHubURL != "" {
		c.GitHub.URL, err = url.Parse(rawGitHubURL)
		if err != nil {
			return
		}
	}

	return
}

func WithConfigFilePath(repository string, path string) ConfigOpt {
	return func() (*Config, error) {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		return WithConfigFileFromReader(repository, f)()
	}
}

var (
	WithDefaultUserConfigFile   = WithConfigFilePath("", DefaultUserConfigFilePath)
	WithDefaultSystemConfigFile = WithConfigFilePath("", DefaultSystemConfigFilePath)
)

func WithConfigFilePaths(repository string, paths ...string) ConfigOpt {
	return func() (*Config, error) {
		opts := []ConfigOpt{}
		for _, path := range paths {
			opts = append(opts, WithConfigFilePath(repository, path))
		}

		return New(opts...)
	}
}

var WithDefaultConfigFiles = WithConfigFilePaths(DefaultUserConfigFilePath, DefaultSystemConfigFilePath)

func WithConfigFileFromReader(repository string, r io.Reader) ConfigOpt {
	return func() (*Config, error) {
		f := &RawConfigFile{}
		if err := toml.NewDecoder(r).Decode(f); err != nil {
			return nil, err
		}

		c, err := f.Parse()
		if err != nil {
			return nil, err
		}

		return c.ToConfig(repository), nil
	}
}

func WithConfigFile(repository string, c *ConfigFile) ConfigOpt {
	return func() (*Config, error) {
		return c.ToConfig(repository), nil
	}
}

func WithConfig(c *Config) ConfigOpt {
	return func() (*Config, error) {
		return c, nil
	}
}
