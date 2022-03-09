package conf

import "net/url"

type Config struct {
	Verbose  bool
	Network  string
	Address  string
	RootDir  string
	StateDir string
	GitHub   *GitHubConfig
	Runtime  *RuntimeConfig
	Secrets  map[string]string
}

type GitHubConfig struct {
	URL   *url.URL
	Token string
}

type RuntimeConfig struct {
	Name  string
	Image string
}
