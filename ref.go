package sequence

import "fmt"

var (
	Name = "sqnc"

	Package = "github.com/frantjc/sequence"

	Registy = "docker.io"

	Repository = "frantjc/sequence"

	Tag = "latest"

	Digest = ""
)

func Image() string {
	ref := Registy
	if Repository != "" {
		ref = fmt.Sprintf("%s/%s", ref, Repository)
	}
	if Tag != "" {
		return fmt.Sprintf("%s:%s", ref, Tag)
	}
	if Digest != "" {
		return fmt.Sprintf("%s@%s", ref, Digest)
	}

	return ref
}
