package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/frantjc/sequence/meta"
)

var (
	conf *Config
)

func Get() (*Config, error) {
	if conf == nil {
		if err := sviper.ReadInConfig(); err == nil {
			ConfigFilePath = sviper.ConfigFileUsed()
		}

		var (
			net          = "unix"
			port         = sviper.GetInt(portKey)
			addr         = sviper.GetString(socketKey)
			runtimeName  = sviper.GetString(runtimeNameKey)
			runtimeImage = sviper.GetString(runtimeImageKey)
			rootDir      = DefaultRootDir
			stateDir     = DefaultStateDir
		)
		if !strings.HasPrefix(ConfigFilePath, sysRootPath) {
			rootDir = filepath.Join(usrRootPath, "lib")
			stateDir = filepath.Join(usrRootPath, "run")
		}
		if port != 0 {
			net = "tcp"
			addr = fmt.Sprintf(":%d", port)
		} else if addr == "" {
			if !strings.HasPrefix(ConfigFilePath, sysRootPath) {
				addr = fmt.Sprintf("unix://%s/%s.sock", stateDir, meta.Name)
			} else {
				addr = DefaultSocket
			}
		}
		if runtimeName == "" {
			runtimeName = DefaultRuntimeName
		}
		if runtimeImage == "" {
			runtimeImage = DefaultRuntimeImage
		}

		conf = &Config{
			Verbose:  sviper.GetBool(verboseKey),
			Network:  net,
			Address:  addr,
			RootDir:  rootDir,
			StateDir: stateDir,
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

		if err := os.MkdirAll(conf.RootDir, 0755); err != nil {
			return nil, err
		}

		if err := os.MkdirAll(conf.StateDir, 0755); err != nil {
			return nil, err
		}
	}

	return conf, nil
}
