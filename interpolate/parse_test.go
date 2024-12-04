package interpolate

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilBindingInterpolate(t *testing.T) {
	result, err := interpolateBytes([]byte("null"), nil)
	assert.NoError(t, err)
	assert.Equal(t, result, nil)
	result, err = interpolateBytes([]byte("true"), nil)
	assert.NoError(t, err)
	assert.Equal(t, result, true)
	result, err = interpolateBytes([]byte("false"), nil)
	assert.NoError(t, err)
	assert.Equal(t, result, false)
	result, err = interpolateBytes([]byte("[]"), nil)
	assert.NoError(t, err)
	assert.Equal(t, result, []interface{}{})
	result, err = interpolateBytes([]byte("{}"), nil)
	assert.NoError(t, err)
	assert.Equal(t, result, map[string]interface{}{})
	result, err = interpolateBytes([]byte("123"), nil)
	assert.NoError(t, err)
	assert.Equal(t, result, float64(123))
}


func TestSimpleBindingInterpolate(t *testing.T) {
	var binding sync.Map
	binding.Store("a", nil)
	binding.Store("b", true)
	binding.Store("c", false)
	binding.Store("d", []interface{}{})
	binding.Store("e", map[string]interface{}{})
	binding.Store("f", float64(123))
	result, err := interpolateBytes([]byte(`"${a}"`), &binding)
	assert.NoError(t, err)
	assert.Equal(t, result, nil)
	result, err = interpolateBytes([]byte(`"${b}"`), &binding)
	assert.NoError(t, err)
	assert.Equal(t, result, true)
	result, err = interpolateBytes([]byte(`"${c}"`), &binding)
	assert.NoError(t, err)
	assert.Equal(t, result, false)
	result, err = interpolateBytes([]byte(`"${d}"`), &binding)
	assert.NoError(t, err)
	assert.Equal(t, result, []interface{}{})
	result, err = interpolateBytes([]byte(`"${e}"`), &binding)
	assert.NoError(t, err)
	assert.Equal(t, result, map[string]interface{}{})
	result, err = interpolateBytes([]byte(`"${f}"`), &binding)
	assert.NoError(t, err)
	assert.Equal(t, result, float64(123))
}

func TestNestedBindingInterpolate(t *testing.T) {
	var binding sync.Map
	binding.Store("e", map[string]interface{}{
		"f": map[string]interface{}{
			"g": nil,
		},
		"h": map[string]interface{}{
			"true": true,
			"false": false,
		},
		"i": []interface{}{},
		"j": map[string]interface{}{},
		"k": 123,
	})
	result, err := interpolateBytes([]byte(`"${e.f.g}"`), &binding)
	assert.NoError(t, err)
	assert.Equal(t, result, nil)
	result, err = interpolateBytes([]byte(`"${e.h.true}"`), &binding)
	assert.NoError(t, err)
	assert.Equal(t, result, true)
	result, err = interpolateBytes([]byte(`"${e.h.false}"`), &binding)
	assert.NoError(t, err)
	assert.Equal(t, result, false)
	result, err = interpolateBytes([]byte(`"${e.i}"`), &binding)
	assert.NoError(t, err)
	assert.Equal(t, result, []interface{}{})
	result, err = interpolateBytes([]byte(`"${e.j}"`), &binding)
	assert.NoError(t, err)
	assert.Equal(t, result, map[string]interface{}{})
	result, err = interpolateBytes([]byte(`"${e.k}"`), &binding)
	assert.NoError(t, err)
	assert.Equal(t, result, float64(123))
}