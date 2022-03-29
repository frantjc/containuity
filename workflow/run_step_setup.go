package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/meta"
	"github.com/frantjc/sequence/runtime"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func runStepSetup(ctx context.Context, r runtime.Runtime, expandedStep *Step, ro *runOpts) (*actions.Metadata, actions.Reference, *runOpts, error) {
	if expandedStep.Uses == "" {
		return nil, nil, ro, nil
	}

	ro.logout.Debugf("[%sSQNC:DBG%s] parsing 'uses: %s'", log.ColorDebug, log.ColorNone, expandedStep.Uses)
	action, err := actions.ParseReference(expandedStep.Uses)
	if err != nil {
		return nil, nil, ro, err
	}

	ro.logout.Infof("[%sSQNC%s] setting up action '%s'", log.ColorInfo, log.ColorNone, action.String())
	spec := &runtime.Spec{
		Image:      meta.Image(),
		Entrypoint: []string{"sqncshim"},
		Cmd:        []string{"plugin", "uses", action.String(), ro.gctx.GitHubContext.ActionPath},
		Mounts: []specs.Mount{
			{
				// actions are global because each step that uses
				// actions/checkout@v2 expects it to function the same
				Source:      getHostActionPath(action, ro),
				Destination: ro.gctx.GitHubContext.ActionPath,
				Type:        runtime.MountTypeBind,
			},
		},
	}
	// make sure all of the host directories that we intend to bind exist
	// note that at this point all bind mounts are directories
	for _, mount := range spec.Mounts {
		if mount.Type == runtime.MountTypeBind {
			if err = os.MkdirAll(mount.Source, 0777); err != nil {
				return nil, nil, ro, err
			}
		}
	}

	var (
		outbuf = new(bytes.Buffer)
		opts   = []runtime.ExecOpt{
			runtime.WithStreams(os.Stdin, outbuf, ro.stderr),
		}
	)
	if err = runSpec(ctx, r, spec, ro, opts); err != nil {
		return nil, nil, ro, err
	}

	resp := &StepOut{}
	if err = json.NewDecoder(outbuf).Decode(resp); err != nil {
		return nil, nil, ro, err
	}

	if actionMetadataJson := []byte(resp.GetActionMetadata()); len(actionMetadataJson) != 0 {
		ro.logout.Debugf("[%sSQNC:DBG%s] parsing metadata for action '%s'", log.ColorDebug, log.ColorNone, action.String())
		actionMetadata := &actions.Metadata{}
		return actionMetadata, action, ro, json.Unmarshal(actionMetadataJson, actionMetadata)
	} else {
		ro.logout.Infof("[%sSQNC:DBG%s] not an action '%s'", log.ColorDebug, log.ColorNone, action.String())
		return nil, nil, ro, actions.ErrNotAnAction
	}
}
