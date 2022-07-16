package plugins

import (
	"io/fs"
	"os"
	"path/filepath"
	"plugin"

	"github.com/frantjc/go-js"
	"github.com/frantjc/sequence/internal/flags"
)

var (
	EnvVarPlugins = "SQNC_PLUGINS"
)

func Open() error {
	var (
		pluginDirs = []string{
			flags.PluginDir,
			"/etc/sqnc/plugins",
			os.Getenv(EnvVarPlugins),
		}
		home, err = os.UserHomeDir()
	)
	if err == nil {
		pluginDirs = append(
			pluginDirs,
			filepath.Join(home, ".sqnc/plugins"),
		)
	}

	for _, dir := range js.Unique(
		js.Filter(pluginDirs, func(s string, _ int, _ []string) bool {
			return s != ""
		}),
	) {
		if pluginDirFi, err := os.Stat(dir); err == nil && pluginDirFi.IsDir() {
			if err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if d.IsDir() || filepath.Ext(d.Name()) != ".so" {
					return nil
				}

				_, e := plugin.Open(path)
				return e
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
