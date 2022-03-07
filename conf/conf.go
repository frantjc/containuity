package conf

import "fmt"

var (
	conf *Config
)

func Get() (*Config, error) {
	if conf == nil {
		err := sviper.ReadInConfig()
		if err != nil {
			return nil, err
		}

		var (
			net          = "unix"
			port         = sviper.GetInt(portKey)
			addr         = sviper.GetString(socketKey)
			runtimeName  = sviper.GetString(runtimeNameKey)
			runtimeImage = sviper.GetString(runtimeImageKey)
		)
		if port != 0 {
			net = "tcp"
			addr = fmt.Sprintf(":%d", port)
		} else if addr == "" {
			addr = DefaultSocket
		}
		if runtimeName == "" {
			runtimeName = DefaultRuntimeName
		}
		if runtimeImage == "" {
			runtimeImage = DefaultRuntimeImage
		}

		conf = &Config{
			Verbose: sviper.GetBool(verboseKey),
			Network: net,
			Address: addr,
			GitHub: &GitHubConfig{
				Token: sviper.GetString(githubTokenKey),
				URL:   DefaultGitHubURL,
			},
			Runtime: &RuntimeConfig{
				Name:  runtimeName,
				Image: runtimeImage,
			},
			Secrets: sviper.GetStringMapString(secretsKey),
		}

		ConfigFilePath = sviper.ConfigFileUsed()
	}

	return conf, nil
}
