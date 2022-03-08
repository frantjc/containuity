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

func JobToTypeJob(j *workflow.Job) *types.Job {
	steps := make([]*types.Step, len(j.Steps))

	for i, s := range j.Steps {
		steps[i] = StepToTypeStep(&s)
	}

	return &types.Job{
		Steps: steps,
	}
}

func WorkflowToTypeWorkflow(w *workflow.Workflow) *types.Workflow {
	jobs := make(map[string]*types.Job, len(w.Jobs))

	for i, j := range w.Jobs {
		jobs[i] = JobToTypeJob(&j)
	}

	return &types.Workflow{
		Jobs: jobs,
	}
}
