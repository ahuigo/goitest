package goitest

import (
	"encoding/json"
	"testing"
)

func TestXxx(t *testing.T) {
	d:=ReqItem{
		// QueryParam: url.Values{
		// 	"key": []string{"value"},
		// },
		Json: map[string]any{
			"name": "Alex",
		},
		Testrules: []TestRule{
		},
	}
	buf, err := json.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("buf:",string(buf))
}