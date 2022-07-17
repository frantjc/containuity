package command_test

import (
	"testing"

	"github.com/frantjc/sequence/pkg/command"
	"github.com/stretchr/testify/assert"
)

func TestNewPruneCommand(t *testing.T) {
	cmd, err := command.NewPruneCmd()
	assert.Nil(t, err)
	assert.NotNil(t, cmd)
}
