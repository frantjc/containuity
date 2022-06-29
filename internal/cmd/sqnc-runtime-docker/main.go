package main

import (
	"github.com/frantjc/sequence/internal/runtimes"
	"github.com/frantjc/sequence/runtime/docker"
)

func init() {
	runtimes.RegisterRuntime(docker.RuntimeName, docker.NewRuntime)
}

func main() {}
