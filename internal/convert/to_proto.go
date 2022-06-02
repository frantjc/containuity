package convert

import (
	"github.com/frantjc/sequence/pb/types"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/workflow"
	"github.com/opencontainers/runtime-spec/specs-go"
	"google.golang.org/protobuf/types/known/anypb"
)

func MapInterfaceToAnyProto(i map[string]interface{}) map[string]*anypb.Any {
	a := map[string]*anypb.Any{}
	for k, v := range i {
		a[k] = v.(*anypb.Any)
	}
	return a
}

func StepToProtoStep(s *workflow.Step) *types.Step {
	if s == nil {
		return nil
	}

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
		Params:     MapInterfaceToAnyProto(s.Params),
	}
}

func JobToProtoJob(j *workflow.Job) *types.Job {
	if j == nil {
		return nil
	}

	steps := make([]*types.Step, len(j.Steps))
	for i, s := range j.Steps {
		steps[i] = StepToProtoStep(s)
	}

	container := &types.Job_Container{}
	if jobImage, ok := j.Container.(string); ok {
		container.Image = jobImage
	}

	return &types.Job{
		Name:      j.Name,
		RunsOn:    j.RunsOn,
		Steps:     steps,
		Outputs:   j.Outputs,
		Env:       j.Env,
		Container: container,
	}
}

func WorkflowToProtoWorkflow(w *workflow.Workflow) *types.Workflow {
	if w == nil {
		return nil
	}

	jobs := make(map[string]*types.Job, len(w.Jobs))
	for i, j := range w.Jobs {
		jobs[i] = JobToProtoJob(j)
	}

	return &types.Workflow{
		Name: w.Name,
		Jobs: jobs,
	}
}

func RuntimeContainerToProtoContainer(container runtime.Container) *types.Container {
	return &types.Container{
		Id: container.ID(),
	}
}

func RuntimeImageToProtoImage(image runtime.Image) *types.Image {
	return &types.Image{
		Ref: image.Ref(),
	}
}

func RuntimeVolumeToProtoVolume(volume runtime.Volume) *types.Volume {
	return &types.Volume{
		Source: volume.Source(),
	}
}

func RuntimeSpecToProtoSpec(s *runtime.Spec) *types.Spec {
	if s == nil {
		return nil
	}

	return &types.Spec{
		Image:      s.Image,
		Cwd:        s.Cwd,
		Entrypoint: s.Entrypoint,
		Cmd:        s.Cmd,
		Env:        s.Env,
		Mounts:     SpecsMountsToProtoMounts(s.Mounts),
		Privileged: s.Privileged,
	}
}

func SpecsMountsToProtoMounts(m []specs.Mount) []*types.Mount {
	mounts := make([]*types.Mount, len(m))
	for i, j := range m {
		mounts[i] = &types.Mount{
			Source:      j.Source,
			Destination: j.Destination,
			Type:        j.Type,
			Options:     j.Options,
		}
	}

	return mounts
}
