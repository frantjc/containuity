package main

import (
	"fmt"
	"io"
	"os"

	"github.com/frantjc/sequence/conf"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

const newline = '\n'

var configGetCmd = &cobra.Command{
	Use:  "get",
	RunE: runConfigGet,
	Args: cobra.RangeArgs(0, 1),
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	c, err := conf.NewFromFlags()
	if err != nil {
		return err
	}

	var (
		key    = ""
		stdout = os.Stdout
	)
	if len(args) > 0 {
		key = args[0]
	}

	switch key {
	case "verbose":
		return write(stdout, c.Verbose)
	case "port":
		return write(stdout, c.Port)
	case "socket":
		return write(stdout, c.Socket)
	case "state_dir":
		return write(stdout, c.StateDir)
	case "root_dir":
		return write(stdout, c.RootDir)
	case "":
		return toml.NewEncoder(stdout).Encode(c.ToConfigFile().Raw())
	default:
		return fmt.Errorf("unrecognized key '%s'", key)
	}
}

func write(w io.Writer, i interface{}) error {
	_, err := w.Write(append([]byte(fmt.Sprint(i)), newline))
	return err
}
