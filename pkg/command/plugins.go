package command

import (
	"io/fs"
	"os"
	"path/filepath"
	"plugin"
)

var (
	EnvVarPlugins = "SQNC_PLUGINS"
)

func OpenPlugins(dir string) error {
	if dirFi, err := os.Stat(dir); err == nil && dirFi.IsDir() {
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

	return nil
}
