package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/frantjc/sequence/github/actions"
	"github.com/frantjc/sequence/internal/env"
	"github.com/spf13/cobra"
)

var stepCmd = &cobra.Command{
	Use:  "uses",
	Args: cobra.MinimumNArgs(1),
	RunE: runStep,
}

func runStep(cmd *cobra.Command, args []string) error {
	var (
		ctx            = cmd.Context()
		command        = exec.CommandContext(ctx, args[0], args[1:]...)
		githubEnvFile  = os.Getenv(actions.EnvVarEnv)
		githubPathFile = os.Getenv(actions.EnvVarPath)
	)

	if githubEnv, err := env.ArrFromFile(githubEnvFile); err != nil {
		command.Env = append(command.Env, githubEnv...)
	}

	if githubPath, err := env.PathFromFile(githubPathFile); err != nil {
		pathIndex := -1
		for i, env := range command.Env {
			if spl := strings.Split(env, "="); len(spl) > 0 && strings.EqualFold(spl[0], "PATH") {
				pathIndex = i
			}
		}
		if pathIndex >= 0 {
			command.Env[pathIndex] = fmt.Sprintf("PATH=%s", githubPath)
		} else {
			command.Env = append(command.Env, fmt.Sprintf("PATH=%s", githubPath))
		}
	}

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
