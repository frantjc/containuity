package sequence

import "path"

const (
	shimName = "shim"
	stateDir = "/run/sqnc"
	shimDir  = stateDir
)

var (
	shimPath        = path.Join(shimDir, shimName)
	actionDir       = path.Join(stateDir, "action")
	actionPath      = actionDir
	workspace       = path.Join(stateDir, "workspace")
	runnerTemp      = path.Join(stateDir, "runner/tmp")
	runnerToolCache = path.Join(stateDir, "runner/toolcache")
)
