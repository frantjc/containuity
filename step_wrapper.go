package sequence

import "github.com/frantjc/sequence/runtime"

type stepWrapper struct {
	step        *Step
	id          string
	extraMounts []*runtime.Mount
	extraEnv    map[string]string
	state       map[string]string
}
