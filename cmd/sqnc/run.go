package main

import (
	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/flags"
	"github.com/frantjc/sequence/internal/plugins"
	"github.com/frantjc/sequence/internal/runtimes"
	"github.com/frantjc/sequence/runtime/docker"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, _ []string) {
			if err := plugins.Open(); err != nil {
				cmd.PrintErrln(err)
			}

			var (
				ctx          = cmd.Context()
				runtime, err = runtimes.GetRuntime(ctx, flags.RuntimeName)
			)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			workflow, err := sequence.NewWorkflowFromFile(flags.File)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			opts := []sequence.ExecutorOpt{sequence.WithRuntime(runtime)}
			if !flags.Quiet {
				opts = append(opts, sequence.WithVerbose)
			}

			executor, err := sequence.NewWorkflowExecutor(ctx, workflow, opts...)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			if err := executor.Execute(ctx); err != nil {
				cmd.PrintErrln(err)
			}
		},
	}
)

func init() {
	runCmd.Flags().StringVarP(&flags.File, "file", "f", "", "file to execute")

	if err := runCmd.MarkFlagFilename("file", "yaml", "yml", "json"); err != nil {
		panic(err)
	}

	if err := runCmd.MarkFlagRequired("file"); err != nil {
		panic(err)
	}
}

func init() {
	runCmd.PersistentFlags().BoolVarP(&flags.Quiet, "quiet", "q", false, "quiet logs")
}

func init() {
	runCmd.PersistentFlags().StringVarP(&flags.RuntimeName, "runtime", "r", docker.RuntimeName, "runtime to use")
}

func init() {
	runCmd.Flags().StringVarP(&flags.PluginDir, "plugins", "p", "", "plugin directory")
	if err := runCmd.MarkFlagDirname("plugins"); err != nil {
		panic(err)
	}
}
