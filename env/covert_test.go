package env_test

import (
	"testing"

	"github.com/frantjc/sequence/env"
	"github.com/stretchr/testify/assert"
)

func TestMapFromArr(t *testing.T) {
	var (
		a        = []string{"KEY1=val", "KEY2=", "=val", "notakeyvalpair"}
		expected = map[string]string{
			"KEY1": "val",
		}
		actual = env.MapFromArr(a)
	)

	assert.Equal(t, expected, actual)
}

func TestToMap(t *testing.T) {
	var (
		expected = map[string]string{
			"KEY1": "val",
		}
		actual = env.ToMap("KEY1=val", "KEY2=", "=val", "notakeyvalpair")
	)

	assert.Equal(t, expected, actual)
}

func TestArrFromMap(t *testing.T) {
	var (
		m = map[string]string{
			"KEY1": "val",
			"":     "val",
			"KEY2": "",
		}
		expected = []string{"KEY1=val"}
		actual   = env.ArrFromMap(m)
	)

	assert.Equal(t, expected, actual)
}

func TestToArr(t *testing.T) {
	var (
		expected = []string{"KEY1=val"}
		actual   = env.ToArr("KEY1", "val", "", "val", "KEY2", "")
	)

	assert.Equal(t, expected, actual)
}
