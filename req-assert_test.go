package goitest

import (
	"bytes"
	"net/http/httptest"
	"os"
	"testing"
)
func TestMain(m *testing.M) {
    code := m.Run() 
	integration.Save()
    os.Exit(code)
}
func TestNewExpectBuilder(t *testing.T) {
	// 1. build request
	req := R(t, "create-user").
		SetJson(map[string]string{"name": "Alex"}).
		SetReq("POST", "http://localhost/api/v1/user")
		//GenCurlCommand()

	// 2. mock response
	req.GenRequest()
	resp := httptest.NewRecorder()
	resp.Body = bytes.NewBuffer([]byte(`{
		"id": 1
	}`))
	resp.Code = 200
	req.SetResponse(resp)

	// 2. assert
	req.AssertBodyContains(`id`)
	req.AssertBodyJqEqual(`.id`, `1`)
	// req.AssertBodyJqEqual(`.id`, `2`)
	req.AssertStatusBetween(200, 300)
	req.Sync()
	req.Test("query-user", func(req *RequestTeser) {
		req.SetQueryParamTpl("id", "create-user.output.id")
		req.SetReq("GET", "http://localhost/api/v1/user")
		req.GenRequest()
		req.Sync()
	})


}
