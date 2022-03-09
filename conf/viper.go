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
	ConfigFileName = "config"

	// ConfigFileName is the type of the file
	// package conf looks for configuration from
	ConfigFileType = "toml"
)

var (
	sviper      = viper.New()
	home        = os.Getenv("HOME")
	usrRootPath = fmt.Sprintf("%s/.%s", home, meta.Name)
	sysRootPath = fmt.Sprintf("/etc/%s", meta.Name)

	// ConfigFilePath is the path to the config file that
	// package conf resolved configuration from
	ConfigFilePath = fmt.Sprintf("%s.%s", ConfigFileName, ConfigFileType)
)

func init() {
	sviper.SetConfigName(ConfigFileName)
	sviper.SetConfigType(ConfigFileType)

	if cwd, err := os.Getwd(); err == nil {
		sviper.AddConfigPath(cwd)
	}
	sviper.AddConfigPath(usrRootPath)
	sviper.AddConfigPath(sysRootPath)

	sviper.SetEnvPrefix(meta.Name)
	sviper.AllowEmptyEnv(true)
}
