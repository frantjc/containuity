package container

import "github.com/opencontainers/runtime-spec/specs-go"

type Mount specs.Mount

const (
	MountOptReadOnly string = "ro"
)

const (
	MountTypeBind string = "bind"
	MountTypeVolume string = "volume"
	MountTypeTmpfs string = "tmpfs"
)
