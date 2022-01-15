package actions_test

import (
	"testing"

	"github.com/frantjc/sequence/github/actions"
	"github.com/stretchr/testify/assert"
)

func TestExpandBytes(t *testing.T) {
	var (
		b       = []byte("${{ github.repository }} ${{ }} ${USER} $USER ${{ }")
		mapping = func(i string) string {
			return "frantjc/sequence"
		}
		expected = []byte("frantjc/sequence  ${USER} $USER ")
		actual   = actions.ExpandBytes(b, mapping)
	)
	assert.Equal(t, expected, actual)
}

func TestExpand(t *testing.T) {
	var (
		s       = "${{ github.repository }} ${{ }} ${HOME} $HOME ${{ }"
		mapping = func(i string) string {
			return "frantjc/sequence"
		}
		expected = "frantjc/sequence  ${HOME} $HOME "
		actual   = actions.Expand(s, mapping)
	)
	assert.Equal(t, expected, actual)
}
