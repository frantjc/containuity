package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/internal/env"
)

func main() {
	if err := mainE(); err != nil {
		panic(err)
	}
}

func mainE() error {
	var (
		ctx            = context.Background()
		args           = os.Args
		command        = exec.CommandContext(ctx, args[0], args[1:]...)
		githubEnvFile  = os.Getenv(actions.EnvVarEnv)
		githubPathFile = os.Getenv(actions.EnvVarPath)
	)

	command.Env = os.Environ()
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if githubEnv, err := env.ArrFromFile(githubEnvFile); err != nil {
		command.Env = append(command.Env, githubEnv...)
	}

	if githubPath, err := env.PathFromFile(githubPathFile); err != nil && githubPath != "" {
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
	}

	return command.Run()
}
