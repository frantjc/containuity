package command

import (
	"github.com/frantjc/sequence/internal/runtimes"
	"github.com/frantjc/sequence/runtime/docker"
	"github.com/spf13/cobra"
)

func NewPruneCmd() (Cmd, error) {
	var (
		runtimeName string
		pruneCmd    = std(&cobra.Command{
			Use:   "prune [--runtime NAME]",
			Short: "Prune dangling resources",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				var (
					ctx = cmd.Context()
				)

				rt, err := runtimes.GetRuntime(ctx, runtimeName)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				if err := rt.PruneContainers(ctx); err != nil {
					cmd.PrintErrln(err)
				}

				if err := rt.PruneVolumes(ctx); err != nil {
					cmd.PrintErrln(err)
				}

				if err := rt.PruneImages(ctx); err != nil {
					cmd.PrintErrln(err)
				}
			},
		})
	)

	flags := pruneCmd.Flags()
	flags.StringVar(&runtimeName, "runtime", docker.RuntimeName, "runtime to use")

	return pruneCmd, nil
}
