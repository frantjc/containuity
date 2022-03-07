package workflow

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/frantjc/sequence/meta"
	"github.com/frantjc/sequence/runtime"
)

var (
	readonly = []string{runtime.MountOptReadOnly}
	workdir  = ""
)

func init() {
	var err error
	workdir, err = os.UserHomeDir()
	if err != nil {
		workdir, _ = os.Getwd()
	}

	workdir = filepath.Join(workdir, fmt.Sprintf(".%s", meta.Name))
}
