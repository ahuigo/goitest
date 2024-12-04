# 🛠️ Go Integration Test
[![tag](https://img.shields.io/github/tag/ahuigo/goitest.svg)](https://github.com/ahuigo/goitest/tags)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-%23007d9c)
[![GoDoc](https://godoc.org/github.com/ahuigo/goitest?status.svg)](https://pkg.go.dev/github.com/ahuigo/goitest)
![Build Status](https://github.com/ahuigo/goitest/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/ahuigo/goitest)](https://goreportcard.com/report/github.com/ahuigo/goitest)
[![Coverage](https://img.shields.io/codecov/c/github/ahuigo/goitest)](https://codecov.io/gh/ahuigo/goitest)
[![Contributors](https://img.shields.io/github/contributors/ahuigo/goitest)](https://github.com/ahuigo/goitest/graphs/contributors)
[![License](https://img.shields.io/github/license/ahuigo/goitest)](./LICENSE)

## Features
- [x] Generate http request in golang
- [x] Generate curl command for http request
- [] Generate integration test cases (and curl)
    - [x] Assert Rule
    - [] Gen integration rules api+data(curl+request data)
            - integration rule struct
            - setParamsFrom: rule.name rule.output(jq)
    - [] integration test server
    - [] integration test ui

## Unittest with gonic

```

func CreateTestCtx(req *http.Request) (resp *httptest.ResponseRecorder, ctx *gin.Context) {
	resp = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(resp)
	ctx.Request = req
	return
}

func TestGonicApi(t *testing.T) {
	// 1. build request
	req, _ := goitest.R().
		SetQueryParams(map[string]string{
			"job_id":   "1234",
		}).
		SetReq("GET", "http://any/api/v1/spark/job").
		GenRequest()
	curl := goitest.GenCurlCommand(req, nil)
	println(curl)
	resp, ctx := CreateTestCtx(req)

	// 2. execute
	sparkServer := GetGonicSparkServer()
	sparkServer.GetJobInfo(ctx)
	if resp.Code != http.StatusOK {
		errors := ctx.Errors.Errors()
		fmt.Println("output", errors)
		t.Errorf("Expect code 200, but get %d body:%v", resp.Code, resp.Body)
	} else {
        data := map[string]string{}
		goitest.BuildResponse(resp.Result()).Json(&data)
		if data["status"] == "" {
			t.Fatalf("Bad response: %v", data)
		}
	}
}
```
# todo

