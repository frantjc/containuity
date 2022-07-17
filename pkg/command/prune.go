package command

import (
	"github.com/frantjc/sequence/internal/runtimes"
	"github.com/frantjc/sequence/runtime/docker"
	"github.com/spf13/cobra"
)

func NewPruneCmd() (Cmd, error) {
	var (
		runtimeName string
		pruneCmd    = &cobra.Command{
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
		}
	)

	flags := pruneCmd.Flags()
	flags.StringVar(&runtimeName, "runtime", docker.RuntimeName, "runtime to use")

	for _, newChildCmd := range []func() (Cmd, error){
		NewPruneContainersCmd,
		NewPruneVolumesCmd,
		NewPruneImagesCmd,
	} {
		childCmd, err := newChildCmd()
		if err != nil {
			return nil, err
		}

		pruneCmd.AddCommand(childCmd.(*cobra.Command))
	}

	return pruneCmd, nil
}

func NewPruneContainersCmd() (Cmd, error) {
	var (
		runtimeName        string
		pruneContainersCmd = &cobra.Command{
			Use:  "containers",
			Args: cobra.NoArgs,
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
			},
		}
	)

	flags := pruneContainersCmd.Flags()
	flags.StringVar(&runtimeName, "runtime", docker.RuntimeName, "runtime to use")

	return pruneContainersCmd, nil
}

func NewPruneVolumesCmd() (Cmd, error) {
	var (
		runtimeName     string
		pruneVolumesCmd = &cobra.Command{
			Use:  "volumes",
			Args: cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				var (
					ctx = cmd.Context()
				)

				rt, err := runtimes.GetRuntime(ctx, runtimeName)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				if err := rt.PruneVolumes(ctx); err != nil {
					cmd.PrintErrln(err)
				}
			},
		}
	)

	flags := pruneVolumesCmd.Flags()
	flags.StringVar(&runtimeName, "runtime", docker.RuntimeName, "runtime to use")

	return pruneVolumesCmd, nil
}

func NewPruneImagesCmd() (Cmd, error) {
	var (
		runtimeName    string
		pruneImagesCmd = &cobra.Command{
			Use:  "images",
			Args: cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				var (
					ctx = cmd.Context()
				)

				rt, err := runtimes.GetRuntime(ctx, runtimeName)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				if err := rt.PruneImages(ctx); err != nil {
					cmd.PrintErrln(err)
				}
			},
		}
	)

	flags := pruneImagesCmd.Flags()
	flags.StringVar(&runtimeName, "runtime", docker.RuntimeName, "runtime to use")

	return pruneImagesCmd, nil
}
