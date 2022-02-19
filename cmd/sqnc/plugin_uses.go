package main

import (
	"encoding/json"
	"os"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/github/actions"
	"github.com/spf13/cobra"
)

var pluginUsesCmd = &cobra.Command{
	RunE: runPluginUses,
	Use:  "uses",
	Args: cobra.RangeArgs(1, 2),
}

func runPluginUses(cmd *cobra.Command, args []string) error {
	var (
		actionRef = args[0]
		path      = "."
	)

	if len(args) > 1 {
		path = args[1]
	}

	parsed, err := actions.ParseReference(actionRef)
	if err != nil {
		return err
	}

	m, err := actions.CloneContext(cmd.Context(), parsed, actions.WithPath(path))
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(&sequence.StepResponse{
		Action: m,
	})
}
