package command

import (
	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/internal/log"
	"github.com/frantjc/sequence/pkg/github/actions"
	"github.com/frantjc/sequence/runtime"
	"github.com/spf13/cobra"
)

func getDefaultExecutorOpts(cmd *cobra.Command, verbose bool) []sequence.ExecutorOpt {
	var (
		stdout = log.New(cmd.OutOrStdout()).SetVerbose(verbose)
		stderr = log.New(cmd.ErrOrStderr()).SetVerbose(verbose)
	)

	return []sequence.ExecutorOpt{
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
}
