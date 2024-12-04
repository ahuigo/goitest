package goitest

import (
	"testing"
)

var (
	jqTestInput = []byte(`{ 
		"foo": [1, 2, 3] ,
		"name": "Alex"
	}
	`)
	jqTestCases = []struct{
		expected string
		expr string
	}{
		{ expected: `[1,2,3]`, expr: ".foo |..", },
		{ expected: `"Alex"`, expr: ".name",},
	}

)
func TestJqRun(t *testing.T) {
	for _, c := range jqTestCases {
		err := jqEqual(c.expr, c.expected, jqTestInput)
		if err != nil {
			t.Fatal(err)
		}
	}

}
