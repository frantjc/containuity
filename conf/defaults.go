package conf

import (
	"fmt"

	"github.com/frantjc/sequence/github"
	"github.com/frantjc/sequence/meta"
)

const (
	DefaultRuntimeName  = "docker"
	DefaultRuntimeImage = "node:12"
	DefaultRootDir      = "/var/lib/sqnc"
	DefaultStateDir     = "/var/run/sqnc"
)

var (
	DefaultSocket = fmt.Sprintf("unix://%s/%s.sock", DefaultStateDir, meta.Name)
)

var (
	DefaultGitHubURL = github.DefaultURL
)
