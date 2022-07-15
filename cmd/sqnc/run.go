package main

import (
	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/runtime/docker"
	"github.com/spf13/cobra"
)

var (
	file    string
	verbose bool
	runCmd  = &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, _ []string) {
			var (
				ctx          = cmd.Context()
				runtime, err = docker.NewRuntime(ctx)
			)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			workflow, err := sequence.NewWorkflowFromFile(file)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			opts := []sequence.ExecutorOpt{sequence.WithRuntime(runtime)}
			if verbose {
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
	runCmd.Flags().StringVarP(&file, "file", "f", "", "Workflow file to execute")
	if err := runCmd.MarkFlagFilename("file", "yaml", "yml", "json"); err != nil {
		panic(err)
	}

	if err := runCmd.MarkFlagRequired("file"); err != nil {
		panic(err)
	}
}

func init() {
	runCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Verbose logs")
}
