package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/env"
	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/internal/shim"
)

func main() {
	if err := mainE(); err != nil {
		os.Exit(1)
	}
}

func mainE() error {
	var (
		ctx  = context.Background()
		args = os.Args
	)

	if runnerToolCache := os.Getenv(actions.EnvVarRunnerToolCache); runnerToolCache != "" {
		os.Setenv("PATH", fmt.Sprintf("%s:%s", runnerToolCache, os.Getenv("PATH")))
	}

	if len(args) == 1 {
		return fmt.Errorf("shim requires at least 1 argument")
	} else if _, ok := os.LookupEnv(shim.EnvVarShimSwitch); !ok {
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

		s, err := json.Marshal(m)
		if err != nil {
			return err
		}

		return json.NewEncoder(os.Stdout).Encode(&sequence.Step_Out{
			Metadata: map[string]string{
				sequence.ActionMetadataKey: string(s),
			},
		})
	}

	var (
		command        = exec.CommandContext(ctx, args[1], args[2:]...) //nolint:gosec
		githubEnvFile  = os.Getenv(actions.EnvVarEnv)
		githubPathFile = os.Getenv(actions.EnvVarPath)
	)

	command.Env = os.Environ()
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if githubEnv, err := env.ArrFromFile(githubEnvFile); err == nil {
		command.Env = append(command.Env, githubEnv...)
	} else {
		if _, err = os.Create(githubEnvFile); err != nil {
			return err
		}
	}

	if githubPath, err := env.PathFromFile(githubPathFile); err == nil && githubPath != "" {
		pathIndex := -1
		for i, env := range command.Env {
			if spl := strings.Split(env, "="); len(spl) > 0 && strings.EqualFold(spl[0], "PATH") {
				pathIndex = i
				break
			}
		}
		if pathIndex >= 0 {
			command.Env[pathIndex] = fmt.Sprintf("PATH=%s:%s", githubPath, os.Getenv("PATH"))
		} else {
			command.Env = append(command.Env, fmt.Sprintf("PATH=%s", githubPath))
		}
	} else {
		if _, err = os.Create(githubPathFile); err != nil {
			return err
		}
	}

	return command.Run()
}
