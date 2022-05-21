package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/frantjc/sequence"
	"github.com/spf13/cobra"
)

const newline = '\n'

var versionCmd = &cobra.Command{
	Use:  "version",
	RunE: runVersion,
}

func runVersion(cmd *cobra.Command, args []string) error {
	return write(os.Stdout, fmt.Sprintf("sqncd%s %s", sequence.Semver(), runtime.Version()))
}

func write(w io.Writer, i interface{}) error {
	_, err := w.Write(append([]byte(fmt.Sprint(i)), newline))
	return err
}
