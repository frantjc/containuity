package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/log"
	"github.com/frantjc/sequence/meta"
	"github.com/frantjc/sequence/runtime"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/spf13/cobra"
)

var runStepCmd = &cobra.Command{
	RunE: runRunStep,
	Use:  "step",
	Args: cobra.MinimumNArgs(1),
}

func init() {
	runStepCmd.Flags().StringVarP(&stepID, "id", "s", "", "ID of the step to run")
	runStepCmd.Flags().StringVarP(&jobName, "job", "j", "", "name of the job to run")
}

func runRunStep(cmd *cobra.Command, args []string) error {
	var (
		ctx  = cmd.Context()
		step *sequence.Step
		job  *sequence.Job
		path = args[0]
		r    io.Reader
		err  error
	)
	if path == fromStdin {
		r = os.Stdin
	} else {
		var err error
		r, err = os.Open(path)
		if err != nil {
			return err
		}
	}

	getConfig()
	if stepID != "" {
		if jobName != "" {
			workflow, err := sequence.NewWorkflowFromReader(r)
			if err != nil {
				return err
			}

			job, err = workflow.GetJob(jobName)
			if err != nil {
				return err
			}
		} else {
			job, err = sequence.NewJobFromReader(r)
			if err != nil {
				return err
			}
		}

		step, err = job.GetStep(stepID)
		if err != nil {
			return err
		}
	} else {
		step, err = sequence.NewStepFromReader(r)
		if err != nil {
			return err
		}
	}

	rt, err := runtime.Get(ctx, runtimeName)
	if err != nil {
		return err
	}

	return runStep(ctx, rt, step, withGitHubToken(gitHubToken))
}

var (
	readonly = []string{runtime.MountOptReadOnly}
	workdir  = ""
)

func init() {
	var err error
	workdir, err = os.UserHomeDir()
	if err != nil {
		workdir, _ = os.Getwd()
	}

	workdir = filepath.Join(workdir, ".sqnc")
}

func runStep(ctx context.Context, r runtime.Runtime, s *sequence.Step, opts ...runOpt) error {
	var (
		ro = &runOpts{
			workflow: &sequence.Workflow{},
			job:      &sequence.Job{},
			path:     ".",
		}
	)
	for _, opt := range opts {
		err := opt(ro)
		if err != nil {
			return err
		}
	}

	vopts := []actions.VarsOpt{actions.WithToken(ro.gitHubToken)}
	ghvars, err := actions.NewVarsFromPath(ro.path, vopts...)
	if err != nil {
		return err
	}

	var (
		ghenv = ghvars.Env
		ghctx = ghvars.ActionsContext
	)
	es, err := expandStep(s.Canonical(), ghctx)
	if err != nil {
		return err
	}

	var (
		// generate a unique, reproducible, directory-name-compliant ID from the current context
		// so that steps that are a part of the same job share the same mounts
		id = base64.URLEncoding.EncodeToString(
			[]byte(ro.path + ro.workflow.Name + ro.jobName),
		)
		gitHubEnv  = filepath.Join(workdir, id, "github", "env")
		gitHubPath = filepath.Join(workdir, id, "github", "path")
		spec       = &runtime.Spec{
			Image:      meta.Image(),
			Cwd:        ghenv.Workspace,
			Privileged: es.Privileged,
			Env:        append(ghenv.Arr(), "SEQUENCE=true"),
			Mounts: []specs.Mount{
				{
					Source:      "/etc/ssl",
					Destination: "/etc/ssl",
					Type:        runtime.MountTypeBind,
					Options:     readonly,
				},
				{
					Source:      filepath.Join(workdir, id, "workspace"),
					Destination: ghenv.Workspace,
					Type:        runtime.MountTypeBind,
				},
				{
					Source:      filepath.Join(workdir, id, "runner", "temp"),
					Destination: ghenv.RunnerTemp,
					Type:        runtime.MountTypeBind,
				},
				{
					Source:      filepath.Join(workdir, id, "runner", "toolcache"),
					Destination: ghenv.RunnerToolCache,
					Type:        runtime.MountTypeBind,
				},
			},
		}
	)

	_, err = createFile(gitHubEnv)
	if err != nil {
		return err
	}

	_, err = createFile(gitHubPath)
	if err != nil {
		return err
	}

	if es.Uses != "" {
		action, err := actions.ParseReference(es.Uses)
		if err != nil {
			return err
		}

		spec.Entrypoint = []string{"sqnc"}
		spec.Cmd = []string{"plugin", "uses", action.String(), ghenv.ActionPath}
		spec.Mounts = append(spec.Mounts, specs.Mount{
			// actions can be global since every step that uses actions/checkout@v2
			// expects to function the same
			Source:      filepath.Join(workdir, "actions", action.Owner(), action.Repository(), action.Path(), action.Version()),
			Destination: ghenv.ActionPath,
			Type:        runtime.MountTypeBind,
		})
	} else {
		if es.Image != "" {
			spec.Image = es.Image
		}

		if es.Run != "" {
			switch es.Shell {
			case "/bin/bash", "bash":
				spec.Entrypoint = []string{"/bin/bash", "-c"}
			case "/bin/sh", "sh", "":
				spec.Entrypoint = []string{"/bin/sh", "-c"}
			default:
				return fmt.Errorf("unsupported shell '%s'", es.Shell)
			}
			spec.Cmd = []string{es.Run}
		} else {
			spec.Entrypoint = es.Entrypoint
			spec.Cmd = es.Cmd
		}
	}

	for _, mount := range spec.Mounts {
		if mount.Type == runtime.MountTypeBind {
			err = os.MkdirAll(mount.Source, 0777)
			if err != nil {
				return err
			}
		}
	}

	// these are _files_, NOT directories
	spec.Mounts = append(spec.Mounts, []specs.Mount{
		{
			Source:      "/etc/hosts",
			Destination: "/etc/hosts",
			Type:        runtime.MountTypeBind,
			Options:     readonly,
		},
		{
			Source:      "/etc/resolv.conf",
			Destination: "/etc/resolv.conf",
			Type:        runtime.MountTypeBind,
			Options:     readonly,
		},
		// these are used like
		// $ echo "MY_VAR=myval" >> $GITHUB_ENV
		// $ echo "/.mybin" >> $GITHUB_PATH
		// respectively. TODO source the contents of these files into spec.Env
		{
			Source:      gitHubEnv,
			Destination: ghenv.Env,
			Type:        runtime.MountTypeBind,
		},
		{
			Source:      gitHubPath,
			Destination: ghenv.Path,
			Type:        runtime.MountTypeBind,
		},
	}...)

	var (
		stdout = log.NewPrefixedInfoWriter(log.ColorInfo + "| " + log.ColorNone)
		popts  = []runtime.PullOpt{runtime.WithStream(stdout)}
		copts  = append([]runtime.SpecOpt{runtime.WithSpec(spec)}, ro.sopts...)
		outbuf = new(bytes.Buffer)
		errbuf = stdout
		eopts  = []runtime.ExecOpt{runtime.WithStreams(os.Stdin, outbuf, errbuf)}
	)
	if !s.IsStdoutResponse() {
		eopts[0] = runtime.WithStreams(os.Stdin, stdout, errbuf)
	}
	log.Infof("%s| running step %s%s", log.ColorInfo, es.GetID(), log.ColorNone)
	if err = runSpec(ctx, r, spec, popts, copts, eopts); err != nil {
		return err
	}

	if s.IsStdoutResponse() {
		resp := &sequence.StepResponse{}
		if err = json.NewDecoder(outbuf).Decode(resp); err != nil {
			return err
		}

		if resp.Action != nil {
			rs, err := sequence.NewStepFromAction(resp.Action, ghenv.ActionPath)
			if err != nil {
				return err
			}

			es, err = expandStep(es.Merge(rs).Canonical(), ghctx)
			if err != nil {
				return err
			}

			spec.Image = es.Image
			spec.Entrypoint = es.Entrypoint
			spec.Cmd = es.Cmd

			for envvar, val := range es.With {
				spec.Env = append(
					spec.Env,
					fmt.Sprintf(
						"INPUT_%s=%s",
						strings.ToUpper(strings.ReplaceAll(envvar, "-", "_")),
						val,
					),
				)
			}

			var (
				copts  = append([]runtime.SpecOpt{runtime.WithSpec(spec)}, ro.sopts...)
				outbuf = stdout
				eopts  = []runtime.ExecOpt{runtime.WithStreams(os.Stdin, outbuf, errbuf)}
			)
			log.Infof("%s| running step %s%s", log.ColorInfo, es.GetID(), log.ColorNone)
			err = runSpec(ctx, r, spec, popts, copts, eopts)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func expandStep(s *sequence.Step, ctx *actions.ActionsContext) (*sequence.Step, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	es := &sequence.Step{}
	err = json.Unmarshal(
		actions.ExpandBytes(b, func(s string) string {
			return ctx.Value(s).(string)
		}),
		es,
	)
	if err != nil {
		return nil, err
	}

	return es, nil
}

func runSpec(ctx context.Context, r runtime.Runtime, s *runtime.Spec, popts []runtime.PullOpt, copts []runtime.SpecOpt, eopts []runtime.ExecOpt) error {
	if s.Image != meta.Image() {
		_, err := r.Pull(ctx, s.Image, popts...)
		if err != nil {
			return err
		}
	}

	container, err := r.Create(ctx, copts...)
	if err != nil {
		return err
	}

	err = container.Exec(ctx, eopts...)
	if err != nil {
		return err
	}

	return nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(name), 0777)
	if err != nil {
		return nil, err
	}

	return os.Create(name)
}
