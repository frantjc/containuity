package conf

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/frantjc/sequence/github"
	"github.com/frantjc/sequence/meta"
	"github.com/frantjc/sequence/runtime"
	"github.com/pelletier/go-toml/v2"
)

//go:embed config.toml
var ExampleRawConfigFileBytes []byte
var ExampleRawConfigFile = &RawConfigFile{}

func init() {
	if err := toml.NewDecoder(bytes.NewReader(ExampleRawConfigFileBytes)).Decode(ExampleRawConfigFile); err != nil {
		panic(fmt.Sprintf("%s/config.ExampleConfigFile is not a valid config file", meta.Module))
	}
}

var DefaultGitHubURL = github.DefaultURL

const (
	DefaultRuntimeName    = runtime.DefaultRuntimeName
	DefaultRunnerImage    = "docker.io/library/node:12"
	DefaultConfigFileName = "config.toml"
)

var home = os.Getenv("HOME")

var (
	DefaultSystemRootDir   = filepath.Join("/var/lib", meta.Name)
	DefaultSystemStateDir  = filepath.Join("/var/run", meta.Name)
	DefaultSystemConfigDir = filepath.Join("/etc", meta.Name)
	DefaultUserDir         = filepath.Join(home, fmt.Sprintf(".%s", meta.Name))
	DefaultUserConfigDir   = DefaultUserDir
	DefaultUserRootDir     = filepath.Join(DefaultUserDir, "lib")
	DefaultUserStateDir    = filepath.Join(DefaultUserDir, "run")
)

var (
	DefaultSystemSocket         = fmt.Sprintf("unix://%s", filepath.Join(DefaultSystemStateDir, fmt.Sprintf("%s.sock", meta.Name)))
	DefaultSystemConfigFilePath = filepath.Join(DefaultSystemConfigDir, DefaultConfigFileName)
	DefaultUserConfigFilePath   = filepath.Join(DefaultUserDir, DefaultConfigFileName)
	DefaultUserSocket           = fmt.Sprintf("unix://%s", filepath.Join(DefaultUserStateDir, fmt.Sprintf("%s.sock", meta.Name)))
)
