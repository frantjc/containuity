package conf

import (
	"fmt"
	"os"

	"github.com/frantjc/sequence/meta"
	"github.com/spf13/viper"
)

const (
	// ConfigFileName is the name of the file
	// package conf looks for configuration from
	ConfigFileName = "config.toml"
)

var (
	sviper = viper.New()

	// ConfigFilePath is the path to the config file that
	// package conf resolved configuration from
	ConfigFilePath = ConfigFileName
)

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		home = "$HOME"
	}

	sviper.SetConfigName(ConfigFileName)
	sviper.SetConfigType("toml")

	if cwd, err := os.Getwd(); err == nil {
		sviper.AddConfigPath(cwd)
	}
	sviper.AddConfigPath(fmt.Sprintf("%s/.%s", home, meta.Name))
	sviper.AddConfigPath(fmt.Sprintf("/etc/%s", meta.Name))

	sviper.SetEnvPrefix(meta.Name)
	sviper.AllowEmptyEnv(true)
}
