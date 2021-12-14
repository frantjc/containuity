package sequence

import "fmt"

var (
	Name = "sqnc"

	Package = "github.com/frantjc/sequence"

	Repository = "frantjc/sequence"

	Tag = "latest"
)

func Image() string {
	if Tag != "" {
		return fmt.Sprintf("%s:%s", Repository, Tag)
	}

	return fmt.Sprintf("%s:latest", Repository)
}
