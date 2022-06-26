package sequence

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/internal/paths"
	"github.com/frantjc/sequence/internal/paths/volumes"
	"github.com/frantjc/sequence/runtime"
)

func (e executor) SetupAction(ctx context.Context, action *actions.Reference) (*actions.Metadata, error) {
	if action == nil {
		return nil, fmt.Errorf("action must not be nil")
	}

	spec := &runtime.Spec{
		Image:      e.RunnerImage.GetRef(),
		Entrypoint: []string{paths.Shim, action.String(), e.GlobalContext.GitHubContext.ActionPath},
		Env: []string{
			"SQNC=true",
			"SEQUENCE=true",
		},
		Mounts: []*runtime.Mount{
			{
				// actions are global because each step that uses
				// actions/checkout@v2 expects it to function the same
				Source:      volumes.GetActionSource(action),
				Destination: e.GlobalContext.GitHubContext.ActionPath,
				Type:        runtime.MountTypeVolume,
			},
		},
	}

	outbuf := new(bytes.Buffer)
	if err := e.RunContainer(ctx, spec, runtime.NewStreams(e.Stdin, outbuf, e.Stderr)); err != nil {
		return nil, err
	}

	out := &Step_Out{}
	if err := json.NewDecoder(outbuf).Decode(out); err != nil {
		return nil, err
	}

	return out.GetActionMetadata()
}
