package e2e_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/meta"
	"github.com/stretchr/testify/assert"
)

var (
	testAction = "frantjc/sequence/testdata@main"
)

func e2eEnabled() bool {
	return os.Getenv("SQNC_E2E") != ""
}

func TestSqncPluginUses(t *testing.T) {
	if e2eEnabled() {
		sqnc := exec.Command(meta.Name, "plugin", "uses", testAction, "/tmp/test/pluginuses")
		buf := new(bytes.Buffer)
		sqnc.Stdout = buf
		err := sqnc.Run()
		assert.Nil(t, err)

		if sqnc.ProcessState != nil {
			assert.True(t, sqnc.ProcessState.Exited())
			assert.Zero(t, sqnc.ProcessState.ExitCode())
		}

		resp := &sequence.StepResponse{}
		err = json.NewDecoder(buf).Decode(resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp.Action)
	}
}
