package command_test

import (
	"testing"

	"github.com/frantjc/sequence/pkg/command"
	"github.com/stretchr/testify/assert"
)

func TestNewAttachCommand(t *testing.T) {
	cmd, err := command.NewAttachCmd()
	assert.Nil(t, err)
	assert.NotNil(t, cmd)
}
