package meta

import (
	"fmt"

	"golang.org/x/mod/semver"
)

var (
	Version = "0.0.0"

	Prerelease = ""

	Build = ""
)

func Semver() string {
	version := Version
	if Prerelease != "" {
		version = fmt.Sprintf("%s-%s", version, Prerelease)
	}
	if Build != "" {
		version = fmt.Sprintf("%s-%s", version, Build)
	}
	return semver.Canonical(version)
}
