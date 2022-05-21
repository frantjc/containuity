package conf

import (
	"fmt"
	"net/url"
	"strings"
)

func New(opts ...ConfigOpt) (*Config, error) {
	c := &Config{}

	for _, opt := range opts {
		o, err := opt()
		if err != nil {
			return nil, err
		}

		c = c.Merge(o)
	}

	return c, nil
}

type ConfigFile struct {
	Verbose      bool
	Port         int64
	Socket       string
	RootDir      string
	StateDir     string
	GitHub       *GitHubConfig
	Runtime      *RuntimeConfig
	Secrets      map[string]string
	Repositories map[string]*ScopedConfig
}

func (c *ConfigFile) ToConfig(repository string) *Config {
	config := &Config{
		Verbose:  c.Verbose,
		Port:     c.Port,
		Socket:   c.Socket,
		RootDir:  c.RootDir,
		StateDir: c.StateDir,
		GitHub:   c.GitHub,
		Runtime:  c.Runtime,
		Secrets:  c.Secrets,
	}

	for k, scopedConfig := range c.Repositories {
		if strings.HasPrefix(repository, k) {
			config = config.Merge(scopedConfig.ToConfig())
		}
	}

	return config
}

func (c *ConfigFile) Raw() *RawConfigFile {
	config := &RawConfigFile{
		Verbose:      c.Verbose,
		Port:         c.Port,
		Socket:       c.Socket,
		RootDir:      c.RootDir,
		StateDir:     c.StateDir,
		GitHub:       c.GitHub.Raw(),
		Runtime:      c.Runtime.Raw(),
		Secrets:      c.Secrets,
		Repositories: make(map[string]*RawScopedConfig, len(c.Repositories)),
	}

	for k, scopedConfig := range c.Repositories {
		config.Repositories[k] = &RawScopedConfig{
			GitHub:  scopedConfig.GitHub.Raw(),
			Runtime: scopedConfig.Runtime.Raw(),
			Secrets: scopedConfig.Secrets,
		}
	}

	return config
}

func (c *ConfigFile) Network() string {
	return network(c.Port, c.Socket)
}

func (c *ConfigFile) Address() string {
	return addr(c.Port, c.Socket)
}

type Config struct {
	Verbose  bool
	Port     int64
	Socket   string
	RootDir  string
	StateDir string
	GitHub   *GitHubConfig
	Runtime  *RuntimeConfig
	Secrets  map[string]string
}

func (c *Config) Network() string {
	return network(c.Port, c.Socket)
}

func (c *Config) Address() string {
	return addr(c.Port, c.Socket)
}

func (c *Config) Merge(config *Config) *Config {
	if !c.Verbose {
		c.Verbose = config.Verbose
	}

	if c.Port == 0 {
		c.Port = config.Port
	}

	if c.Socket == "" {
		c.Socket = config.Socket
	}

	if c.RootDir == "" {
		c.RootDir = config.RootDir
	}

	if c.StateDir == "" {
		c.StateDir = config.StateDir
	}

	if c.GitHub == nil {
		c.GitHub = config.GitHub
	} else {
		if c.GitHub.Token == "" {
			c.GitHub.Token = config.GitHub.Token
		}

		if c.GitHub.URL == nil {
			c.GitHub.URL = config.GitHub.URL
		}
	}

	if c.Runtime == nil {
		c.Runtime = config.Runtime
	} else {
		if c.Runtime.Name == "" {
			c.Runtime.Name = config.Runtime.Name
		}

		if c.Runtime.RunnerImage == "" {
			c.Runtime.RunnerImage = config.Runtime.RunnerImage
		}
	}

	if c.Secrets == nil || len(c.Secrets) == 0 {
		c.Secrets = config.Secrets
	} else {
		for k, v := range config.Secrets {
			if _, ok := c.Secrets[k]; !ok {
				c.Secrets[k] = v
			}
		}
	}

	return c
}

func (c *Config) ToConfigFile() *ConfigFile {
	return &ConfigFile{
		Verbose:  c.Verbose,
		Port:     c.Port,
		Socket:   c.Socket,
		RootDir:  c.RootDir,
		StateDir: c.StateDir,
		GitHub:   c.GitHub,
		Runtime:  c.Runtime,
		Secrets:  c.Secrets,
	}
}

type ScopedConfig struct {
	GitHub  *GitHubConfig
	Runtime *RuntimeConfig
	Secrets map[string]string
}

func (c *ScopedConfig) ToConfig() *Config {
	return &Config{
		GitHub:  c.GitHub,
		Runtime: c.Runtime,
		Secrets: c.Secrets,
	}
}

type GitHubConfig struct {
	URL   *url.URL
	Token string
}

func (c *GitHubConfig) Raw() *RawGitHubConfig {
	return &RawGitHubConfig{
		URL:   c.URL.String(),
		Token: c.Token,
	}
}

type RuntimeConfig struct {
	Name        string
	RunnerImage string
}

func (c *RuntimeConfig) Raw() *RawRuntimeConfig {
	return &RawRuntimeConfig{
		Name:        c.Name,
		RunnerImage: c.RunnerImage,
	}
}

const (
	tcp  = "tcp"
	unix = "unix"
)

func network(port int64, socket string) string {
	if port < 65535 && port > 0 {
		return tcp
	}

	return unix
}

func addr(port int64, socket string) string {
	if port < 65535 && port > 0 {
		return fmt.Sprintf(":%d", port)
	}

	return socket
}
