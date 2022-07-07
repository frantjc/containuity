package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/frantjc/go-js"
	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/shim"
	"github.com/frantjc/sequence/pkg/envconv"
	"github.com/frantjc/sequence/pkg/github/actions"
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

	if len(args) == 1 {
		return fmt.Errorf("%s requires at least 1 argument", os.Args[0])
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

	if githubEnv, err := envconv.ArrFromFile(githubEnvFile); err == nil {
		command.Env = append(command.Env, githubEnv...)
	} else {
		if _, err = os.Create(githubEnvFile); err != nil {
			return err
		}
	}

	if githubPath, err := envconv.PathFromFile(githubPathFile); err == nil && githubPath != "" {
		pathIndex := js.FindIndex(command.Env, func(s string, _ int, _ []string) bool {
			spl := strings.Split(s, "=")
			return len(spl) > 0 && strings.EqualFold(spl[0], "PATH")
		})

		pathAddendum := ""

		if runnerToolCache := os.Getenv(actions.EnvVarRunnerToolCache); runnerToolCache != "" {
			pathAddendum = fmt.Sprintf(":%s", runnerToolCache)
		}

		if pathIndex >= 0 {
			command.Env[pathIndex] = fmt.Sprintf("PATH=%s:%s%s", githubPath, os.Getenv("PATH"), pathAddendum)
		} else {
			command.Env = append(command.Env, fmt.Sprintf("PATH=%s%s", githubPath, pathAddendum))
		}
	} else {
		if _, err = os.Create(githubPathFile); err != nil {
			return err
		}
	}

	return command.Run()
}
