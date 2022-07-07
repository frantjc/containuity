package envconv_test

import (
	"testing"

	"github.com/frantjc/go-js"
	"github.com/frantjc/sequence/pkg/envconv"
	"github.com/stretchr/testify/assert"
)

func TestMapFromArr(t *testing.T) {
	var (
		a        = []string{"KEY1=val", "KEY2=", "=val", "notakeyvalpair"}
		expected = map[string]string{
			"KEY1": "val",
			"KEY2": "",
		}
		actual = envconv.MapFromArr(a)
	)

	assert.Equal(t, expected, actual)
}

func TestToMap(t *testing.T) {
	var (
		expected = map[string]string{
			"KEY1": "val",
			"KEY2": "",
		}
		actual = envconv.ToMap("KEY1=val", "KEY2=", "=val", "notakeyvalpair")
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
		expected = []string{"KEY1=val", "KEY2="}
		actual   = envconv.ArrFromMap(m)
	)

	assert.True(t, js.Every(expected, func(e string, _ int, _ []string) bool {
		return js.Includes(actual, e)
	}))
}

func TestToArr(t *testing.T) {
	var (
		expected = []string{"KEY2=", "KEY1=val"}
		actual   = envconv.ToArr("KEY1", "val", "", "val", "KEY2", "")
	)

	assert.True(t, js.Every(expected, func(e string, _ int, _ []string) bool {
		return js.Includes(actual, e)
	}))
}
