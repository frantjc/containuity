package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var runPluginUsesCmd = &cobra.Command{
	RunE: runRunPluginUses,
	SilenceErrors:    true,
	SilenceUsage:     true,
	Use:              "uses",
}

func runRunPluginUses(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		var (
			uses = args[0]
			path = ""
		)

		if len(args) > 1 {
			path = args[1]
		}

		opts := []string{"clone", "--single-branch"}
		usesSpl := strings.Split(uses, "@")
		action := usesSpl[0]
		if len(usesSpl) > 1 {
			revision := usesSpl[1]
			opts = append(opts, "--branch", revision)
		}
		opts = append(opts, fmt.Sprintf("https://github.com/%s", action))
		if path != "" && path != "." && path != "./" {
			opts = append(opts, path)
		}

		cmd := exec.Command("git", opts...)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return nil
}
