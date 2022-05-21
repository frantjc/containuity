package conf

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/frantjc/sequence/github"
	"github.com/frantjc/sequence/runtime"
	"github.com/pelletier/go-toml/v2"
)

var (
	//go:embed config
	ExampleRawConfigFileBytes []byte
	ExampleRawConfigFile      = &RawConfigFile{}
)

func init() {
	if err := toml.NewDecoder(bytes.NewReader(ExampleRawConfigFileBytes)).Decode(ExampleRawConfigFile); err != nil {
		panic("github.com/frantjc/sequence/internal/conf.ExampleConfigFile is not a valid config file")
	}
}

var DefaultGitHubURL = github.DefaultURL

const (
	DefaultRuntimeName    = runtime.DefaultRuntimeName
	DefaultRunnerImage    = "docker.io/library/node:12"
	DefaultConfigFileName = "config"
)

var (
	home, _ = os.UserHomeDir()
	name    = "sqnc"
)

var (
	DefaultSystemRootDir   = filepath.Join("/var/lib", name)
	DefaultSystemStateDir  = filepath.Join("/var/run", name)
	DefaultSystemConfigDir = filepath.Join("/etc", name)
	DefaultUserDir         = filepath.Join(home, fmt.Sprintf(".%s", name))
	DefaultUserConfigDir   = DefaultUserDir
	DefaultUserRootDir     = filepath.Join(DefaultUserDir, "lib")
	DefaultUserStateDir    = filepath.Join(DefaultUserDir, "run")
)

var (
	DefaultSystemSocket         = fmt.Sprintf("unix://%s", filepath.Join(DefaultSystemStateDir, fmt.Sprintf("%s.sock", name)))
	DefaultSystemConfigFilePath = filepath.Join(DefaultSystemConfigDir, DefaultConfigFileName)
	DefaultUserSocket           = fmt.Sprintf("unix://%s", filepath.Join(DefaultUserStateDir, fmt.Sprintf("%s.sock", name)))
	DefaultUserConfigFilePath   = filepath.Join(DefaultUserDir, DefaultConfigFileName)
)
