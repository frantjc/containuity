package sequence

import "fmt"

var (
	Name = "sqnc"

	Module = "github.com/frantjc/sequence"

	Registry = "docker.io"

	Repository = "frantjc/sequence"

	Tag = "latest"

	Digest = ""
)

func init() {
	if Registry == "" {
		panic(fmt.Sprintf("%s.Registry must not be empty", Module))
	}
	if Repository == "" {
		panic(fmt.Sprintf("%s.Repository must not be empty", Module))
	}
}

func Image() string {
	ref := fmt.Sprintf("%s/%s", Registry, Repository)
	if Tag != "" {
		return fmt.Sprintf("%s:%s", ref, Tag)
	}
	if Digest != "" {
		return fmt.Sprintf("%s@%s", ref, Digest)
	}

	return ref
}
