package main

import (
	"fmt"
	"os"
	"os/exec"

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

	if githubPath, err := env.PathFromFile(githubPathFile); err != nil {
		command.Path = fmt.Sprintf("%s:%s", command.Path, githubPath)
	}

	if githubEnv, err := env.ArrFromFile(githubEnvFile); err != nil {
		command.Env = append(command.Env, githubEnv...)
	}

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
