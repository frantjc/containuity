package main

import (
	"encoding/json"
	"os"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/github/actions"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var pluginUsesCmd = &cobra.Command{
	RunE:    runPluginUses,
	Use:     "uses",
	Args:    cobra.RangeArgs(1, 2),
	PreRunE: preRunPluginUses,
}

func preRunPluginUses(cmd *cobra.Command, args []string) error {
	log.Logger = log.Output(os.Stderr)
	return nil
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
		log.Debug().Err(err).Msg("parsing action failed")
		return err
	}

	log.Debug().Msgf("cloning %s to %s", parsed.String(), path)
	a, err := actions.CloneContext(cmd.Context(), parsed, actions.WithPath(path))
	if err != nil {
		log.Debug().Err(err).Msg("clone failed")
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(&sequence.StepResponse{
		Action: a,
	})
}
