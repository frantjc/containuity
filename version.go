package sequence

import (
	"fmt"
)

var (
	Name = "sqnc"

	Package = "github.com/frantjc/sequence"

	SemVer = "0.0.0"

	Revision = ""

	Version = ""
)

func init() {
	v := SemVer
	if Revision != "" {
		v = fmt.Sprintf("%s+%s", v, Revision)
	}
	Version = v
}
