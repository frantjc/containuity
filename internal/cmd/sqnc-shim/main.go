package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/frantjc/go-js"
	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/pkg/envconv"
	"github.com/frantjc/sequence/pkg/github/actions"
	"github.com/frantjc/sequence/pkg/github/actions/uses"
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

	if len(args) < 2 {
		return fmt.Errorf("'%s' requires at least 2 arguments, e.g. '%s -e echo hello there'", args[0], args[0])
	}

	switch args[1] {
	// clone
	case "-c":
		if len(args) < 3 {
			return fmt.Errorf("'%s %s' requires at least 1 argument, e.g. '%s %s actions/checkout@v2'", args[0], args[1], args[0], args[1])
		}

		var (
			usesStr = args[2]
			path    = "."
		)

		if len(args) > 3 {
			path = args[3]
		}

		parsed, err := uses.Parse(usesStr)
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
	// hang
	case "-s":
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		close(sigs)
		return nil
	// exec
	case "-e":
		if len(args) < 3 {
			return fmt.Errorf("'%s %s' requires at least 1 argument, e.g. '%s %s echo hello there'", args[0], args[1], args[0], args[1])
		}

		var (
			command        = exec.CommandContext(ctx, args[2], args[3:]...) //nolint:gosec
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

	return fmt.Errorf("unrecognized argument '%s'", args[1])
}
