package workflow

import "github.com/frantjc/sequence/conf"

type workflowOpts struct {
	conf *conf.Config
}

type WorkflowOpt func(wo *workflowOpts) error

func WithConfig(c *conf.Config) WorkflowOpt {
	return func(wo *workflowOpts) error {
		wo.conf = c
		return nil
	}
}
