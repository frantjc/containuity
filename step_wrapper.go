package sequence

import "github.com/frantjc/sequence/runtime"

type StepWrapper struct {
	*Step
	ExtraMounts []*runtime.Mount
	State       map[string]string
}
