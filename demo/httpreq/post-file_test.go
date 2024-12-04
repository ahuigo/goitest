package goitest

import (
	"strings"
	"testing"

	goitest "github.com/ahuigo/goitest"
)

func TestPostFile(t *testing.T) {
	curl, err := goitest.R(t, "post_file").
		SetQueryParams(map[string]string{"p": "1"}).
		SetFormData(map[string]string{"key": "xx"}).
		SetAuthBasic("user", "pass").
		SetHeader("header1", "value1").
		AddCookieKV("count", "1").
		AddFileHeader("file", "test.txt", []byte("hello world")).
		AddFile("file2", getTestDataPath("a.txt")).
		SetReq("GET", "/path").
		GenCurlCommand()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(curl, "curl ") {
		t.Fatal("bad curl: ", curl)
	} else {
		t.Log(curl)
	}
}
