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

//nolint: gocyclo
func NewRunCommand() (Cmd, error) {
	var (
		workflowFile string
		runtimeName  string
		verbose      bool
		githubToken  string
		context      string
		runCmd       = &cobra.Command{
			Use:  "run",
			Args: cobra.NoArgs,
			Run: func(cmd *cobra.Command, _ []string) {
				var (
					ctx    = cmd.Context()
					stdout = log.New(cmd.OutOrStdout())
					stderr = log.New(cmd.ErrOrStderr())
				)

				rt, err := runtimes.GetRuntime(ctx, runtimeName)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				workflow, err := sequence.NewWorkflowFromFile(workflowFile)
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

				var (
					echo bool
					opts = []sequence.ExecutorOpt{
						sequence.WithRuntime(rt),
						sequence.WithGlobalContext(gc),
						sequence.WithStreams(
							cmd.InOrStdin(),
							cmd.OutOrStdout(),
							cmd.ErrOrStderr(),
						),
						sequence.OnImagePull(func(i runtime.Image) {
							stdout.Infof("[%sSQNC:INFO%s] pulling image '%s'", log.ColorInfo, log.ColorNone, i.GetRef())
						}),
						sequence.OnWorkflowCommand(func(wc *actions.WorkflowCommand) {
							switch wc.Command {
							case actions.CommandError:
								stderr.Infof("[%sACTN:ERR%s] %s", log.ColorError, log.ColorNone, wc.Value)
							case actions.CommandWarning:
								stderr.Infof("[%sACTN:WRN%s] %s", log.ColorWarn, log.ColorNone, wc.Value)
							case actions.CommandNotice:
								stdout.Infof("[%sACTN:NTC%s] %s", log.ColorNotice, log.ColorNone, wc.Value)
							case actions.CommandDebug:
								if verbose || echo {
									stdout.Infof("[%sACTN:DBG%s] %s", log.ColorDebug, log.ColorNone, wc.Value)
								}
							case actions.CommandSetOutput:
								if verbose || echo {
									stdout.Infof("[%sSQNC:DBG%s] %s %s=%s for", log.ColorDebug, log.ColorNone, wc.Command, wc.GetName(), wc.Value)
								}
							case actions.CommandStopCommands:
								if verbose || echo {
									stdout.Infof("[%sSQNC:DBG%s] %s until '%s'", log.ColorDebug, log.ColorNone, wc.Command, wc.Value)
								}
							case actions.CommandEcho:
								if wc.Value == "on" {
									echo = true
								} else if wc.Value == "off" {
									echo = false
								}
							case actions.CommandSaveState:
								if verbose || echo {
									stdout.Infof("[%sSQNC:DBG%s] %s %s=%s", log.ColorDebug, log.ColorNone, wc.Command, wc.GetName(), wc.Value)
								}
							default:
								if verbose || echo {
									stdout.Infof("[%sSQNC:DBG%s] swallowing unrecognized workflow command '%s'", log.ColorDebug, log.ColorNone, wc.Command)
								}
							}
						}),
					}
				)
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

	runCmd.Flags().StringVarP(&workflowFile, "file", "f", "", "workflow file to execute")
	if err := runCmd.MarkFlagFilename("file", "yaml", "yml", "json"); err != nil {
		return nil, err
	}
	if err := runCmd.MarkFlagRequired("file"); err != nil {
		return nil, err
	}

	runCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "debug logs")
	runCmd.PersistentFlags().StringVar(&runtimeName, "runtime", docker.RuntimeName, "runtime to use")
	runCmd.PersistentFlags().StringVar(&githubToken, "github-token", "", "GitHub token to use")
	runCmd.PersistentFlags().StringVar(&context, "context", "", "path to get context from .git")

	return runCmd, nil
}
