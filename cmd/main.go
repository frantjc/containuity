package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/frantjc/containuity"
)

var cmd containuity.ContainuityCmd

func init() {
	cmd.Version = containuity.V
}

func main() {
	var cmd containuity.ContainuityCmd
	cmd.Version = containuity.V

	parser := flags.NewParser(&cmd, flags.HelpFlag)
	parser.NamespaceDelimiter = "-"
	_, err := parser.Parse() 
	if err == flags.ErrHelp {
		fmt.Println(err)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
