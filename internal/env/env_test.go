package env_test

import (
	"testing"

	"github.com/frantjc/sequence/internal/env"
	"github.com/stretchr/testify/assert"
)

func TestArrToMap(t *testing.T) {
	var (
		a = []string{"KEY=val", "KEY=", "=val", "notakeyvalpair"}
		expected = map[string]string{
			"KEY": "val",
		}
		actual = env.ArrToMap(a)
	)
	
	assert.Equal(t, expected, actual)
}

func TestToMap(t *testing.T) {
	var (
		expected = map[string]string{
			"KEY": "val",
		}
		actual = env.ToMap("KEY=val", "KEY=", "=val", "notakeyvalpair")
	)
	
	assert.Equal(t, expected, actual)
}

func TestMapToArr(t *testing.T) {
	var (
		m = map[string]string{
			"KEY1": "val",
			"": "val",
			"KEY2": "",
		}
		expected = []string{"KEY1=val"}
		actual = env.MapToArr(m)
	)
	
	assert.Equal(t, expected, actual)
}
