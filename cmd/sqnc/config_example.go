package main

import (
	"os"

	"github.com/frantjc/sequence/conf"
	"github.com/spf13/cobra"
)

var configExampleCmd = &cobra.Command{
	Use:  "example",
	RunE: runConfigExample,
}

func runConfigExample(cmd *cobra.Command, args []string) error {
	_, err := os.Stdout.Write(conf.ExampleRawConfigFileBytes)
	return err
}
