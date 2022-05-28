package conf

import "net/url"

type RawConfigFile struct {
	Verbose      bool                        `toml:"verbose,omitempty"`
	Port         int64                       `toml:"port,omitempty"`
	HTTPPort     int64                       `toml:"http_port,omitempty"`
	Socket       string                      `toml:"socket,omitempty"`
	RootDir      string                      `toml:"root_dir,omitempty"`
	StateDir     string                      `toml:"state_dir,omitempty"`
	GitHub       *RawGitHubConfig            `toml:"github,omitempty"`
	Runtime      *RawRuntimeConfig           `toml:"runtime,omitempty"`
	Secrets      map[string]string           `toml:"secrets,omitempty"`
	Repositories map[string]*RawScopedConfig `toml:"repository,omitempty"`
}

func (r *RawConfigFile) Parse() (*ConfigFile, error) {
	if r.GitHub == nil {
		r.GitHub = &RawGitHubConfig{}
	}
	github, err := r.GitHub.Parse()
	if err != nil {
		return nil, err
	}

	if r.Runtime == nil {
		r.Runtime = &RawRuntimeConfig{}
	}
	runtime, err := r.Runtime.Parse()
	if err != nil {
		return nil, err
	}

	if r.Port > 65535 || r.Port < 0 {
		return nil, ErrInvalidPort
	}

	if r.HTTPPort > 65535 || r.HTTPPort < 0 {
		return nil, ErrInvalidPort
	}

	repositories := map[string]*ScopedConfig{}
	for scope, rawScopedConfig := range r.Repositories {
		repositories[scope], err = rawScopedConfig.Parse()
		if err != nil {
			return nil, err
		}
	}

	return &ConfigFile{
		Verbose:      r.Verbose,
		Port:         r.Port,
		Socket:       r.Socket,
		RootDir:      r.RootDir,
		StateDir:     r.StateDir,
		GitHub:       github,
		Runtime:      runtime,
		Secrets:      r.Secrets,
		Repositories: repositories,
	}, nil
}

type RawScopedConfig struct {
	GitHub  *RawGitHubConfig  `toml:"github,omitempty"`
	Runtime *RawRuntimeConfig `toml:"runtime,omitempty"`
	Secrets map[string]string `toml:"secrets,omitempty"`
}

func (r *RawScopedConfig) Parse() (*ScopedConfig, error) {
	if r.GitHub == nil {
		r.GitHub = &RawGitHubConfig{}
	}
	github, err := r.GitHub.Parse()
	if err != nil {
		return nil, err
	}

	if r.Runtime == nil {
		r.Runtime = &RawRuntimeConfig{}
	}
	runtime, err := r.Runtime.Parse()
	if err != nil {
		return nil, err
	}

	return &ScopedConfig{
		GitHub:  github,
		Runtime: runtime,
		Secrets: r.Secrets,
	}, nil
}

type RawGitHubConfig struct {
	URL   string `toml:"base_url,omitempty"`
	Token string `toml:"token,omitempty"`
}

func (r *RawGitHubConfig) Parse() (*GitHubConfig, error) {
	parsedURL, err := url.Parse(r.URL)
	if err != nil {
		return nil, err
	}

	return &GitHubConfig{
		URL:   parsedURL,
		Token: r.Token,
	}, nil
}

type RawRuntimeConfig struct {
	Name        string `toml:"name,omitempty"`
	RunnerImage string `toml:"runner_image,omitempty"`
}

func (r *RawRuntimeConfig) Parse() (*RuntimeConfig, error) {
	return &RuntimeConfig{
		Name:        r.Name,
		RunnerImage: r.RunnerImage,
	}, nil
}
