package main

import (
	"fmt"
	"os"

	"github.com/frantjc/sequence/pkg/command"
)

func main() {
	rootCmd, err := command.NewRootCmd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
