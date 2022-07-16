package main

import (
	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/flags"
	"github.com/frantjc/sequence/internal/runtimes"
	"github.com/frantjc/sequence/runtime/docker"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, _ []string) {
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
			if flags.Verbose {
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
	runCmd.Flags().StringVarP(&flags.File, "file", "f", "", "workflow file to execute")

	if err := runCmd.MarkFlagFilename("file", "yaml", "yml", "json"); err != nil {
		panic(err)
	}

	if err := runCmd.MarkFlagRequired("file"); err != nil {
		panic(err)
	}
}

func init() {
	runCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "V", false, "debug logs")
}

func init() {
	runCmd.PersistentFlags().StringVar(&flags.RuntimeName, "runtime", docker.RuntimeName, "runtime to use")
}
