package containuity

import (
	"fmt"
	"os"
)

var Version = "0.0.0"

var Commit = ""

var GoVersion = ""

// V prints the version of containuity embedded within the binary at build time and exits.
func V() {
	v := fmt.Sprintf("containuity %s", Version)
	if Commit != "" {
		v = fmt.Sprintf("%s %s", v, Commit)
	}
	if GoVersion != "" {
		v = fmt.Sprintf("%s %s", v, GoVersion)
	}
	fmt.Println(v)
	os.Exit(0)
}
