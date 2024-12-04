package interpolate

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var environment = map[string]interface{} {
	"a": map[string]interface{}{
		"b": map[string]interface{}{
			"c": []interface{}{1.0, 2.0, 3.0},
			"d": "simple string",
			"e": true,
		},
		"f": false,
		"g": nil,
	},
	"h": 1.23,
}

var testValidExpressions = []struct{
	Source string
	Expect interface{}
}{
	{"noop", "noop"},
	{"${a}", environment["a"]},
	{"${a.b}", environment["a"].(map[string]interface{})["b"]},
	{"${a.b.c}", environment["a"].(map[string]interface{})["b"].(map[string]interface{})["c"]},
	{"${a.b.d}", environment["a"].(map[string]interface{})["b"].(map[string]interface{})["d"]},
	{"${a.b.e}", environment["a"].(map[string]interface{})["b"].(map[string]interface{})["e"]},
	{"${a.f}", environment["a"].(map[string]interface{})["f"]},
	{"${a.g}", environment["a"].(map[string]interface{})["g"]},
	{"${h}", environment["h"]},
	{"${a.g}&${h}", "&1.23"},
	{"${h}", 1.23},
	{"${if h>1 then 100 else 101 end}", "${if h>1 then 100 else 101 end}"},
	{"${a.g}&${h}&${a.f}", "&1.23&false"},
	{"${a.g}&${h}&${a.f}&${a.b.e}&${a.b.d}", "&1.23&false&true&simple string"},
}

var testInvalidExpressions = []struct{
	Source string
	Expect interface{}
} {
	{"${noop}", nil},
	{"${noop}&${no}", "&"},
}

func TestInterpolation(t *testing.T) {
	envbytes, _ := json.Marshal(environment)
	for _, expr := range testValidExpressions {
		result, err := Interpolation(expr.Source, envbytes)
		assert.NoError(t, err)
		assert.Equal(t, expr.Expect, result)
	}
	for _, expr := range testInvalidExpressions {
		result, err := Interpolation(expr.Source, envbytes)
		assert.NoError(t, err)
		assert.Equal(t, expr.Expect, result)
	}
}