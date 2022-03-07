package actions

import (
	"errors"
	"strings"
)

const (
	CommandDebug        = "debug"
	CommandGroup        = "group"
	CommandEndGroup     = "endgroup"
	CommandAddMatcher   = "addmatcher"
	CommandSaveState    = "save-state"
	CommandSetOutput    = "set-output"
	CommandNotice       = "notice"
	CommandWarning      = "warning"
	CommandError        = "error"
	CommandAddMask      = "add-mask"
	CommandEcho         = "echo"
	CommandStopCommands = "stop-commands"
)

var (
	ErrNotACommand = errors.New("not a workflow command")
)

func ErrIsNotACommand(err error) bool {
	return errors.Is(err, ErrNotACommand)
}

func ParseStringCommand(s string) (*Command, error) {
	if !strings.HasPrefix(s, "::") {
		return nil, ErrNotACommand
	}

	a := strings.Split(s, "::")
	if len(a) < 2 {
		return nil, ErrNotACommand
	}

	cmdAndParams := a[1]
	b := strings.Split(cmdAndParams, " ")
	if len(b) < 1 {
		return nil, ErrNotACommand
	}

	cmd := b[0]
	params := map[string]string{}

	if len(b) > 1 {
		for _, p := range strings.Split(b[1], ",") {
			if f := strings.Split(p, "="); len(f) > 0 {
				params[f[0]] = f[1]
			}
		}
	}

	value := ""
	if len(a) > 2 {
		value = a[2]
	}

	return &Command{
		Command:    cmd,
		Parameters: params,
		Value:      value,
	}, nil
}

func ParseCommand(b []byte) (*Command, error) {
	return ParseStringCommand(string(b))
}
