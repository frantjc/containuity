package main

import (
	"fmt"

	"github.com/frantjc/sequence/internal/actions"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var runPluginUsesCmd = &cobra.Command{
	RunE: runRunPluginUses,
	Use:  "uses",
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

		parsed, err := actions.Parse(uses)
		if err != nil {
			return err
		}

		a, err := actions.CloneContext(cmd.Context(), parsed, actions.WithPath(path))
		if err != nil {
			log.Err(err).Msg("")
			return err
		}

		fmt.Println(a)
	}

	return nil
}
