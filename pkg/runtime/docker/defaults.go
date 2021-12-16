package runtime

import (
	dmount "github.com/docker/docker/api/types/mount"
	"github.com/frantjc/sequence"
)

func defaultLabels() map[string]string {
	return map[string]string{
		sequence.Name: "true",
	}
}

func defaultMounts() []dmount.Mount {
	return []dmount.Mount{
		{
			Type:   dmount.TypeTmpfs,
			Target: "/dev/shm",
			TmpfsOptions: &dmount.TmpfsOptions{
				Mode: 01777,
			},
		},
	}
}
