package sequence

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/frantjc/sequence/internal/paths"
	"github.com/frantjc/sequence/internal/paths/volumes"
	"github.com/frantjc/sequence/pkg/github/actions"
	"github.com/frantjc/sequence/pkg/github/actions/uses"
	"github.com/frantjc/sequence/runtime"
)

func (e *executor) SetupAction(ctx context.Context, action *uses.Uses) (*actions.Metadata, error) {
	if action == nil {
		return nil, fmt.Errorf("action must not be nil")
	}

	var (
		spec = &runtime.Spec{
			Image:      e.RunnerImage.GetRef(),
			Entrypoint: []string{paths.Shim, "-c", action.String(), e.GlobalContext.GitHubContext.ActionPath},
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
		outbuf = new(bytes.Buffer)
		out    = &Step_Out{}
	)
	if err := e.RunContainer(ctx, spec, runtime.NewStreams(e.StreamIn, outbuf, e.StreamErr)); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(outbuf).Decode(out); err != nil {
		return nil, err
	}

	return out.GetActionMetadata()
}
