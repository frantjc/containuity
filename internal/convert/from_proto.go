package convert

import (
	"github.com/frantjc/sequence/api/types"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/workflow"
	"github.com/opencontainers/runtime-spec/specs-go"
)

// func MapAnyToInterface(a map[string]*anypb.Any) map[string]interface{}

func ProtoStepToStep(s *types.Step) *workflow.Step {
	return &workflow.Step{
		ID:         s.Id,
		Name:       s.Name,
		Env:        s.Env,
		Shell:      s.Shell,
		Run:        s.Run,
		Uses:       s.Uses,
		With:       s.With,
		Image:      s.Image,
		Entrypoint: s.Entrypoint,
		Cmd:        s.Cmd,
		Privileged: s.Privileged,
		Get:        s.Get,
		Put:        s.Put,
		// TODO cast this somehow
		// Params:     s.Params,
	}
}

func ProtoJobToJob(j *types.Job) *workflow.Job {
	steps := make([]*workflow.Step, len(j.Steps))

	for i, s := range j.Steps {
		steps[i] = ProtoStepToStep(s)
	}

	return &workflow.Job{
		Name:    j.Name,
		RunsOn:  j.RunsOn,
		Steps:   steps,
		Outputs: j.Outputs,
		Env:     j.Env,
		Container: &struct{ Image string }{
			Image: j.Container.Image,
		},
	}
}

func ProtoWorkflowToWorkflow(w *types.Workflow) *workflow.Workflow {
	jobs := make(map[string]*workflow.Job, len(w.Jobs))

	for i, j := range w.Jobs {
		jobs[i] = ProtoJobToJob(j)
	}

	return &workflow.Workflow{
		Name: w.Name,
		Jobs: jobs,
	}
}

func ProtoSpecToRuntimeSpec(s *types.Spec) *runtime.Spec {
	return &runtime.Spec{
		Image:      s.Image,
		Cwd:        s.Cwd,
		Entrypoint: s.Entrypoint,
		Cmd:        s.Cmd,
		Env:        s.Env,
		Mounts:     ProtoMountsToSpecsMounts(s.Mounts),
		Privileged: s.Privileged,
	}
}

func ProtoMountsToSpecsMounts(m []*types.Mount) []specs.Mount {
	mounts := make([]specs.Mount, len(m))

	for i, j := range m {
		mounts[i] = specs.Mount{
			Source:      j.Source,
			Destination: j.Destination,
			Type:        j.Type,
			Options:     j.Options,
		}
	}

	return mounts
}
