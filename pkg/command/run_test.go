package command_test

import (
	"testing"

	"github.com/frantjc/sequence/pkg/command"
	"github.com/stretchr/testify/assert"
)

func TestNewRunCommand(t *testing.T) {
	cmd, err := command.NewRunCommand()
	assert.Nil(t, err)
	assert.NotNil(t, cmd)
}
