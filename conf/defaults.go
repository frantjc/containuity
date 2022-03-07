package conf

import (
	"fmt"
	"os"

	"github.com/frantjc/sequence/github"
	"github.com/frantjc/sequence/meta"
)

const (
	DefaultRuntimeName  = "docker"
	DefaultRuntimeImage = "node:12"
)

var (
	DefaultSocket = fmt.Sprintf("unix://%s/.%s/%s.sock", os.Getenv("HOME"), meta.Name, meta.Name)
)

var (
	DefaultGitHubURL = github.DefaultURL
)
