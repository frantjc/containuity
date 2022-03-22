package main

import (
	"fmt"
	"io"
	"runtime"

	"github.com/frantjc/sequence/meta"
	"github.com/spf13/cobra"
)

const newline = '\n'

var versionCmd = &cobra.Command{
	Use:  "version",
	RunE: runVersion,
}

func runVersion(cmd *cobra.Command, args []string) error {
	return write(cmd.OutOrStdout(), fmt.Sprintf("%s%s %s", meta.Name, meta.Semver(), runtime.Version()))
}

func write(w io.Writer, i interface{}) error {
	_, err := w.Write(append([]byte(fmt.Sprint(i)), newline))
	return err
}