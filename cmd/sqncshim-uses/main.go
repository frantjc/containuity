package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/workflow"
)

func main() {
	if err := mainE(); err != nil {
		panic(err)
	}
}

func mainE() error {
	var (
		args      = os.Args
		ctx       = context.Background()
		actionRef = args[0]
		path      = "."
	)

	if len(args) > 1 {
		path = args[1]
	}

	parsed, err := actions.ParseReference(actionRef)
	if err != nil {
		return err
	}

	m, err := actions.CloneContext(ctx, parsed, actions.WithPath(path))
	if err != nil {
		return err
	}

	s, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(&workflow.StepOut{
		Metadata: map[string]string{
			workflow.ActionMetadataKey: string(s),
		},
	})
}
