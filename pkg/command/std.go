package command

import (
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
)

func std(cmd *cobra.Command) *cobra.Command {
	cmd.SetOut(colorable.NewColorableStdout())
	cmd.SetErr(colorable.NewColorableStderr())
	return cmd
}
