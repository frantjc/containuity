package sequence

import (
	"fmt"
)

var (
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
