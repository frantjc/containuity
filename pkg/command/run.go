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

func NewRunCmd() (Cmd, error) {
	var (
		workflowFile string
		runtimeName  string
		verbose      bool
		githubToken  string
		context      string
		runCmd       = &cobra.Command{
			Use:   "run -f WORKFLOW_FILE [-V] [--github-token STRING] [--context DIR] [--runtime NAME]",
			Short: "Run a workflow file",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, _ []string) {
				var (
					ctx    = cmd.Context()
					stdout = log.New(cmd.OutOrStdout()).SetVerbose(verbose)
					stderr = log.New(cmd.ErrOrStderr()).SetVerbose(verbose)
				)

				rt, err := runtimes.GetRuntime(ctx, runtimeName)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				var workflow *sequence.Workflow
				if workflowFile == "-" {
					workflow, err = sequence.NewWorkflowFromReader(os.Stdin)
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

				opts := []sequence.ExecutorOpt{
					sequence.WithRuntime(rt),
					sequence.WithGlobalContext(gc),
					sequence.WithStreams(
						cmd.InOrStdin(),
						cmd.OutOrStdout(),
						cmd.ErrOrStderr(),
					),
					sequence.OnImagePull(func(event *sequence.Event[runtime.Image]) {
						stdout.Infof("[%sSQNC:INF%s] pulling image '%s'", log.ColorInfo, log.ColorNone, event.Type.GetRef())
					}),
					sequence.OnStepStart(func(event *sequence.Event[*sequence.Step]) {
						stdout.Infof("[%sSQNC:INF%s] running step '%s'", log.ColorInfo, log.ColorNone, event.Type.GetID())
					}),
					sequence.OnJobStart(func(event *sequence.Event[*sequence.Job]) {
						stdout.Infof("[%sSQNC:INF%s] running job '%s'", log.ColorInfo, log.ColorNone, event.GlobalContext.GitHubContext.Job)
					}),
					sequence.OnWorkflowStart(func(event *sequence.Event[*sequence.Workflow]) {
						stdout.Infof("[%sSQNC:INF%s] running workflow '%s'", log.ColorInfo, log.ColorNone, event.GlobalContext.GitHubContext.Workflow)
					}),
					sequence.OnWorkflowCommand(func(event *sequence.Event[*actions.WorkflowCommand]) {
						switch event.Type.Command {
						case actions.CommandError:
							stderr.Infof("[%sACTN:ERR%s] %s", log.ColorError, log.ColorNone, event.Type.Value)
						case actions.CommandWarning:
							stderr.Infof("[%sACTN:WRN%s] %s", log.ColorWarn, log.ColorNone, event.Type.Value)
						case actions.CommandNotice:
							stdout.Infof("[%sACTN:NTC%s] %s", log.ColorNotice, log.ColorNone, event.Type.Value)
						case actions.CommandDebug:
							stdout.Debugf("[%sACTN:DBG%s] %s", log.ColorDebug, log.ColorNone, event.Type.Value)
						case actions.CommandSetOutput:
							stdout.Debugf("[%sSQNC:DBG%s] %s %s=%s for", log.ColorDebug, log.ColorNone, event.Type.Command, event.Type.GetName(), event.Type.Value)
						case actions.CommandStopCommands:
							stdout.Debugf("[%sSQNC:DBG%s] %s until '%s'", log.ColorDebug, log.ColorNone, event.Type.Command, event.Type.Value)
						case actions.CommandEcho:
							switch event.Type.Value {
							case "on":
								stdout.SetVerbose(true)
								stderr.SetVerbose(true)
							case "off":
								stdout.SetVerbose(false)
								stderr.SetVerbose(false)
							default:
								stdout.Debugf("[%sSQNC:DBG%s] swallowing unrecognized value '%s' for workflow command '%s', must be 'on' or 'off'", log.ColorDebug, log.ColorNone, event.Type.Value, event.Type.Command)
							}
						case actions.CommandSaveState:
							stdout.Debugf("[%sSQNC:DBG%s] %s %s=%s", log.ColorDebug, log.ColorNone, event.Type.Command, event.Type.GetName(), event.Type.Value)
						default:
							stdout.Debugf("[%sSQNC:DBG%s] swallowing unrecognized workflow command '%s'", log.ColorDebug, log.ColorNone, event.Type.Command)
						}
					}),
				}
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
		}
	)

	flags := runCmd.Flags()
	flags.StringVarP(&workflowFile, "file", "f", "", "workflow file to execute")
	if err := runCmd.MarkFlagFilename("file", "yaml", "yml", "json"); err != nil {
		return nil, err
	}
	if err := runCmd.MarkFlagRequired("file"); err != nil {
		return nil, err
	}

	persistentFlags := runCmd.PersistentFlags()
	persistentFlags.BoolVarP(&verbose, "verbose", "V", false, "debug logs")
	persistentFlags.StringVar(&runtimeName, "runtime", docker.RuntimeName, "runtime to use")
	persistentFlags.StringVar(&githubToken, "github-token", "", "GitHub token to use")
	persistentFlags.StringVar(&context, "context", "", "path to get context from .git")

	return runCmd, nil
}
