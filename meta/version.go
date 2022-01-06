package meta

import (
	"fmt"

	"golang.org/x/mod/semver"
)

var (
	Major = 0

	Minor = 0

	Patch = 0

	Prerelease = ""

	Build = ""
)

func Version() string {
	version := fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)
	if Prerelease != "" {
		version = fmt.Sprintf("%s-%s", version, Prerelease)
	}
	if Build != "" {
		version = fmt.Sprintf("%s-%s", version, Build)
	}
	return semver.Canonical(version)
}
