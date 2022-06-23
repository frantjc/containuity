package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/github/actions"
)

func main() {
	if err := mainE(); err != nil {
		panic(err)
	}
}

func mainE() error {
	var (
		ctx  = context.Background()
		args = os.Args
	)

	if len(args) == 1 {
		return fmt.Errorf("sqncshim requires at least 1 argument")
	}

	var (
		actionRef = args[1]
		path      = "."
	)

	if len(args) > 1 {
		path = args[2]
	}

	parsed, err := actions.ParseReference(actionRef)
	if err != nil {
		return err
	}

	m, err := actions.CloneContext(ctx, parsed, actions.WithPath(path))
	if err != nil {
		return err
	}

	// for some reason, 's' is flagged as being unused
	// despite it being used below ~line 57
	s, err := json.Marshal(m) //nolint:typecheck
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(&sequence.Step_Out{
		Metadata: map[string]string{
			sequence.ActionMetadataKey: string(s),
		},
	})
}
