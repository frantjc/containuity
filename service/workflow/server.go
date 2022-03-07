package workflow

import (
	"context"

	jobapi "github.com/frantjc/sequence/api/v1/job"
	stepapi "github.com/frantjc/sequence/api/v1/step"
	workflowapi "github.com/frantjc/sequence/api/v1/workflow"
	"github.com/frantjc/sequence/conf"
	"github.com/frantjc/sequence/runtime"
	"google.golang.org/grpc"
)

func Register(ctx context.Context, s grpc.ServiceRegistrar, opts ...WorkflowOpt) error {
	var (
		so  = &workflowOpts{}
		err error
	)
	for _, opt := range opts {
		err := opt(so)
		if err != nil {
			return err
		}
	}

	if so.conf == nil {
		so.conf, err = conf.Get()
		if err != nil {
			return err
		}
	}

	r, err := runtime.Get(ctx, so.conf.Runtime.Name)
	if err != nil {
		return err
	}

	w := &workflowServer{
		conf: so.conf,
		r:    r,
	}

	workflowapi.RegisterWorkflowServer(s, w)
	stepapi.RegisterStepServer(s, w)
	jobapi.RegisterJobServer(s, w)

	return nil
}

type workflowServer struct {
	jobapi.UnimplementedJobServer
	stepapi.UnimplementedStepServer
	workflowapi.UnimplementedWorkflowServer
	conf *conf.Config
	r    runtime.Runtime
}

var _ jobapi.JobServer = &workflowServer{}
var _ stepapi.StepServer = &workflowServer{}
var _ workflowapi.WorkflowServer = &workflowServer{}
