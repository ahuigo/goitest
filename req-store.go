package goitest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/ahuigo/goitest/interpolate"
	"github.com/emirpasic/gods/v2/maps/linkedhashmap"
)

type ReqItem struct {
	Name   string
	Url    string
	Method string
	// QueryParam  url.Values
	// FormData    url.Values
	IsMultiPart   bool
	InputHeaders  map[string]string
	InputParams   map[string]string
	InputFormData map[string]string
	InputBody     string
	Json          any
	OutputHeaders map[string]string
	OutputBody    string
	// request 模板
	Tpl       *RequestTpl
	Testrules []TestRule
}

const (
	integrationTestDataPath = "./tmp/integration-data.json"
)

type _IntegrationTest struct {
	fh       *os.File
	reqitems *linkedhashmap.Map[string, ReqItem]
	reqs     *linkedhashmap.Map[string, *RequestTeser]
}

var (
	integration          *_IntegrationTest
	_integrationTestOnce sync.Once
)

func getIntegrationTest() *_IntegrationTest {
	_integrationTestOnce.Do(func() {
		// fh, err := os.OpenFile(integrationTestDataPath, os.O_WRONLY|os.O_CREATE, 0644)
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Println("data1:",string(content))

		// defer fh.Close()
		integration = &_IntegrationTest{
			// fh:    fh,
			reqitems: linkedhashmap.New[string, ReqItem](),
			reqs:     linkedhashmap.New[string, *RequestTeser](),
		}
		content, err := os.ReadFile(integrationTestDataPath)
		if err == nil && len(content) > 0 {
			json.Unmarshal(content, &integration.reqitems)
			integration.initReqCase()
		}
	})
	return integration
}

func (rt *_IntegrationTest) initReqCase() {
	for _, item := range integration.reqitems.Values() {
		integration.reqs.Put(item.Name, &RequestTeser{
			name: item.Name,
			url:  item.Url,
			rawreq: &http.Request{
				Method:     item.Method,
				Header:     make(http.Header),
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
			},
		})
	}
}

func (rt *_IntegrationTest) AddReqCase(req *RequestTeser) {
	rt.reqs.Put(req.name, req)
	headers := map[string]string{}
	inputParams := map[string]string{}
	inputFormData := map[string]string{}
	for k, vs := range req.rawreq.URL.Query() {
		for _, v := range vs {
			inputParams[k] = v
		}
	}
	for k, vs := range req.rawreq.Header {
		for _, v := range vs {
			headers[k] = v
		}
	}
	for k, vs := range req.formData {
		for _, v := range vs {
			inputFormData[k] = v
		}
	}
	outputBody := ""
	outputHeaders := map[string]string{}
	if req.resp != nil {
		for k, vs := range req.resp.Header() {
			for _, v := range vs {
				outputHeaders[k] = v
			}
		}
		outputBody = req.resp.Body.String()
	}

	rt.reqitems.Put(req.name, ReqItem{
		Name:          req.name,
		Url:           req.url,
		Method:        req.rawreq.Method,
		IsMultiPart:   req.isMultiPart,
		Json:          req.json,
		InputHeaders:  headers,
		InputParams:   inputParams,
		InputFormData: inputFormData,
		// InputBody:     nil,
		OutputHeaders: outputHeaders,
		OutputBody:    outputBody,
		Tpl:           req.tpl,
		Testrules:     req.testrules,
	})
}
func (rt *_IntegrationTest) GetValue(expr string) (any, error) {
		segs := strings.SplitN(expr, ".", 3)
		if len(segs) != 3 {
			return nil, errors.New("Invalid jq expr: " + expr)
		}
		name, part, jqExpr:=segs[0], segs[1],segs[2]
		req, found := rt.reqitems.Get(name)
		if !found {
			return nil, errors.New("Request not found: " + name)
		}
		switch part {
		case "output":
				a:=[]byte(req.OutputBody)
                expr := fmt.Sprintf("${%s}", jqExpr)
				return interpolate.Interpolation(expr, a)

		}
		return nil, errors.New("Invalid jq expr: " + expr)



}

func (rt *_IntegrationTest) Save() {
	// for _, req := range rt.reqs.Values() {
	// }
	values:= rt.reqitems.Values()
	buf, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		panic(err)
	}
	if rt.fh == nil {
		rt.fh, err = os.OpenFile(integrationTestDataPath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
	}
	_, err = rt.fh.Write(buf)
	if err != nil {
		panic(err)
	}
}
