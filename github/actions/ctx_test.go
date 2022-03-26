package actions_test

import (
	"context"
	"testing"

	"github.com/frantjc/sequence/github/actions"
	"github.com/stretchr/testify/assert"
)

func TestWithContext(t *testing.T) {
	var (
		ctx      = context.Background()
		expected = &actions.GlobalContext{
			JobContext: &actions.JobContext{
				Status: "test",
			},
		}
		actual, err = actions.Context(actions.WithContext(ctx, expected))
	)
	assert.Nil(t, err)
	assert.Equal(t, expected.JobContext.Status, actual.JobContext.Status)
}
