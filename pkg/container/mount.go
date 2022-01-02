package container

import "github.com/opencontainers/runtime-spec/specs-go"

type Mount specs.Mount

const (
	MountOptReadOnly = "ro"
)

const (
	MountTypeBind   = "bind"
	MountTypeVolume = "volume"
	MountTypeTmpfs  = "tmpfs"
)
