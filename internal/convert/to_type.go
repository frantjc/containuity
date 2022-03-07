package convert

import (
	"github.com/frantjc/sequence/api/types"
	"github.com/frantjc/sequence/workflow"
)

func StepToTypeStep(s *workflow.Step) *types.Step {
	return &types.Step{
		Id:         s.ID,
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
		Params:     MapInterfaceToAny(s.Params),
	}
}
