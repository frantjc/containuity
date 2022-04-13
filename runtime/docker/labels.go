package docker

import "github.com/docker/docker/api/types/filters"

var (
	labels = map[string]string{
		"sequence": "true",
	}
	filter = filters.NewArgs(
		filters.Arg("label", "sequence=true"),
	)
)
