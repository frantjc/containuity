package plugin

import "errors"

var (
	ErrNoType          = errors.New("plugin: no type")
	ErrNoPluginID      = errors.New("plugin: no id")
	ErrIDRegistered    = errors.New("plugin: id already registered")
	ErrSkipPlugin      = errors.New("skip plugin")
	ErrInvalidRequires = errors.New("invalid requires")
)

func ErrIsSkipPlugin(err error) bool {
	return errors.Is(err, ErrSkipPlugin)
}
