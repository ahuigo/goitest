package goitest

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	goitest "github.com/ahuigo/goitest"
)

func TestResponse(t *testing.T) {
	response := http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		Body:       io.NopCloser(bytes.NewBuffer([]byte(`{"name":"ahuigo"}`))),
	}
	user := struct {
		Name string `json:"name"`
	}{}
	res := goitest.BuildResponse(&response)
	res.Json(&user)
	if user.Name != "ahuigo" {
		t.Fatalf("json parse error: %v", user)
	}
}
