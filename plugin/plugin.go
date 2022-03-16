package plugin

import (
	"context"
	"fmt"
	"path/filepath"
	"plugin"
)

type Type string

func (t Type) String() string { return string(t) }

const (
	Runtime Type = "runtime"
	Service Type = "service"
)

type Registration struct {
	Type      Type
	ID        string
	DependsOn []Type
	InitF     func(context.Context) error
}

func (r *Registration) Init() error {
	return nil
}

func (r *Registration) String() string {
	return fmt.Sprintf("%s.%s", r.Type, r.ID)
}

func Load(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	pattern := filepath.Join(abs, "*")

	libs, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, lib := range libs {
		if _, err := plugin.Open(lib); err != nil {
			return err
		}
	}

	return nil
}
