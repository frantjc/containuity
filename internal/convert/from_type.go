package convert

import (
	"github.com/frantjc/sequence/api/types"
	"github.com/frantjc/sequence/workflow"
	"google.golang.org/protobuf/types/known/anypb"
)

func MapInterfaceToAny(i map[string]interface{}) map[string]*anypb.Any {
	a := map[string]*anypb.Any{}
	for k, v := range i {
		a[k] = v.(*anypb.Any)
	}
	return a
}

// func MapAnyToInterface(a map[string]*anypb.Any) map[string]interface{}

func TypeStepToStep(s *types.Step) *workflow.Step {
	return &workflow.Step{
		ID:         s.Id,
		Name:       s.Name,
		Image:      s.Image,
		Entrypoint: s.Entrypoint,
		Cmd:        s.Cmd,
		Privileged: s.Privileged,
		Env:        s.Env,
		Shell:      s.Shell,
		Run:        s.Run,
		Uses:       s.Uses,
		With:       s.With,
		Get:        s.Get,
		Put:        s.Put,
		// TODO cast this somehow
		// Params:     s.Params,
	}
}
