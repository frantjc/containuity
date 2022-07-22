package command

import (
	"os"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/log"
	"github.com/frantjc/sequence/internal/runtimes"
	"github.com/frantjc/sequence/pkg/github/actions"
	"github.com/frantjc/sequence/runtime"
	"github.com/frantjc/sequence/runtime/docker"
	"github.com/spf13/cobra"
)

func NewAttachCmd() (Cmd, error) {
	var (
		workflowFile string
		runtimeName  string
		verbose      bool
		githubToken  string
		context      string
		runCmd       = std(&cobra.Command{
			Use:   "attach -f WORKFLOW_FILE [-V] [--github-token STRING] [--context DIR] [--runtime NAME]",
			Short: "Attach to a workflow file",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, _ []string) {
				var (
					ctx    = cmd.Context()
					stdout = log.New(cmd.OutOrStdout()).SetVerbose(verbose)
				)

				rt, err := runtimes.GetRuntime(ctx, runtimeName)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				var workflow *sequence.Workflow
				if workflowFile == "-" {
					workflow, err = sequence.NewWorkflowFromReader(cmd.InOrStdin())
				} else {
					workflow, err = sequence.NewWorkflowFromFile(workflowFile)
				}
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				if githubToken == "" {
					githubToken = os.Getenv("GITHUB_TOKEN")
				}

				if context == "" {
					context, err = os.Getwd()
					if err != nil {
						cmd.PrintErrln(err)
						return
					}
				}

				gc, err := actions.NewContextFromPath(ctx, context, actions.WithToken(githubToken))
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				opts := append(
					getDefaultExecutorOpts(cmd, verbose),
					sequence.WithRuntime(rt),
					sequence.WithGlobalContext(gc),
					sequence.OnContainerCreate(func(event *sequence.Event[runtime.Container]) {
						stdout.Infof("[%sSQNC:INF%s] attaching to step", log.ColorError, log.ColorNone)
						if err := event.Type.Attach(ctx, &runtime.Streams{
							In:  cmd.InOrStdin(),
							Out: cmd.OutOrStdout(),
							Err: cmd.ErrOrStderr(),
						}); err != nil {
							cmd.PrintErrln(err)
						}
						stdout.Infof("[%sSQNC:INF%s] detached from step", log.ColorError, log.ColorNone)
					}),
				)
				if verbose {
					opts = append(opts, sequence.WithVerbose)
				}

				executor, err := sequence.NewWorkflowExecutor(ctx, workflow, opts...)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				if err := executor.ExecuteContext(ctx); err != nil {
					cmd.PrintErrln(err)
				}
			},
		})
	)

	flags := runCmd.Flags()
	flags.BoolVarP(&verbose, "verbose", "V", false, "debug logs")
	flags.StringVar(&runtimeName, "runtime", docker.RuntimeName, "runtime to use")
	flags.StringVar(&githubToken, "github-token", "", "GitHub token to use")
	flags.StringVar(&context, "context", "", "path to get context from .git")
	flags.StringVarP(&workflowFile, "file", "f", "", "workflow file to execute")
	if err := runCmd.MarkFlagFilename("file", "yaml", "yml", "json"); err != nil {
		return nil, err
	}
	if err := runCmd.MarkFlagRequired("file"); err != nil {
		return nil, err
	}

	return runCmd, nil
}
