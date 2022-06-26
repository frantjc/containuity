package paths

import (
	"path"

	"github.com/frantjc/sequence/internal/shim"
)

const (
	Root = "/run/sqnc"
)

var (
	Shim            = path.Join(Root, shim.Name)
	Action          = path.Join(Root, "action")
	Workspace       = path.Join(Root, "workspace")
	RunnerTemp      = path.Join(Root, "runner/tmp")
	RunnerToolCache = path.Join(Root, "runner/toolcache")
)
