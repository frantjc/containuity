package defaults

import "github.com/opencontainers/runtime-spec/specs-go"

var (
	Labels = map[string]string{
		"SEQUENCE": "true",
	}
	Mounts = []specs.Mount{
		// {
		// 	Destination: "/dev/shm",
		// 	Type:        runtime.MountTypeTmpfs,
		// 	Options:     []string{"nosuid", "noexec", "nodev", "mode=1777"},
		// },
	}
)
