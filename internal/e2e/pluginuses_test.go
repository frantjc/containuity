package e2e_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/frantjc/sequence"
	"github.com/stretchr/testify/assert"
)

var (
	TestAction = "frantjc/sequence/internal/testdata@main"
)

func init() {
	passedTestAction := os.Getenv("SQNC_TEST_ACTION")
	if passedTestAction != "" {
		TestAction = passedTestAction
	}
}

func TestSqncPluginUses(t *testing.T) {
	sqnc := exec.Command(sequence.Name, "plugin", "uses", TestAction, "/tmp/test/pluginuses")
	buf := new(bytes.Buffer)
	sqnc.Stdout = buf
	err := sqnc.Run()
	assert.Nil(t, err)

	assert.True(t, sqnc.ProcessState.Exited())
	assert.Zero(t, sqnc.ProcessState.ExitCode())

	resp := &sequence.StepResponse{}
	err = json.NewDecoder(buf).Decode(resp)
	assert.Nil(t, err)
	assert.NotNil(t, resp.Step)
}
