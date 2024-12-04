package goitest

import (
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ExpectTestType string

const (
	ExpectTestHeaderEqual    ExpectTestType = "header-equal"
	ExpectTestHeaderContains ExpectTestType = "header-contains"
	ExpectTestStatusBetween  ExpectTestType = "status-between"
	ExpectTestBodyEqual      ExpectTestType = "body-equal"
	ExpectTestBodyContains   ExpectTestType = "body-contains"
	ExpectTestBodyJqEqual    ExpectTestType = "body-jq-equal"
	ExpectTestBodyJqContains ExpectTestType = "body-jq-contains"
)

type TestRule struct {
	AssertType         ExpectTestType
	PropOrExpr         string
	ExpectedVal        string
	ExpectedBetweenInt [2]int
}

func (r *RequestTeser) SetResponse(resp *httptest.ResponseRecorder) *RequestTeser {
	if resp == nil {
		r.t.Error("ResponseRecorder is nil")
	}
	r.Response = resp.Result()
	r.resp = resp
	return r
}
func (r *RequestTeser) CreateGinContext() (ctx *gin.Context) {
    r.t.Helper()
	req, err := r.GenRequest()
    if err != nil {
        r.t.Fatal(err)
    }
	resp := httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(resp)
	ctx.Request = req
	return ctx
}



type HeaderType string

const (
	HeaderContentType   HeaderType = "Content-Type"
	HeaderAuthorization HeaderType = "Authorization"
	HeaderCookie        HeaderType = "Cookie"
)

func (r *RequestTeser) AssertHeaderEqual(prop HeaderType, val string) *RequestTeser {
	if r.Response.Header.Get(string(prop)) != val {
		err := errors.Errorf("expect header %s: %s, got %s", prop, val, r.Response.Header.Get(string(prop)))
		r.t.Fatal(err)
	}
	return r
}


// ExpectHeaderContains("Cookie", "session=123")
func (r *RequestTeser) AssertHeaderContains(prop HeaderType, substr string) *RequestTeser {
	val := r.Response.Header.Get(string(prop))
	if strings.Contains(val, substr) {
		err := errors.Errorf("expect header %s: %s, got %s", prop, substr, r.Response.Header.Get(string(prop)))
		r.t.Fatal(err)
	}
	return r
}

func (r *RequestTeser) AssertStatusBetween(start, end int) *RequestTeser {
	status := r.Response.StatusCode
	if status < start || status > end {
		err := errors.Errorf("expect status between %d and %d, got %d", start, end, status)
		r.t.Fatal(err)
	}
	return r
}
func (r *RequestTeser) AssertBodyContains(substr string) *RequestTeser {
	r.t.Helper()
	r.testrules = append(r.testrules, TestRule{
		AssertType:  ExpectTestBodyContains,
		PropOrExpr:  substr,
	})
	if !strings.Contains(r.resp.Body.String(), substr) {
		err :=errors.Errorf("expect body contains %s, got %s", substr, r.resp.Body.String())
		r.t.Fatal(err)
	}
	return r
}
func (r *RequestTeser) AssertBodyJqEqual(expr string, expectedVal string) *RequestTeser {
	r.t.Helper()
	r.testrules = append(r.testrules, TestRule{
		AssertType:  ExpectTestBodyJqEqual,
		PropOrExpr:  expr,
		ExpectedVal: expectedVal,
	})
	err := jqEqual(expr, expectedVal, r.resp.Body.Bytes())
	if err != nil {
		r.t.Fatal(err)
	}
	return r
}

func (r *RequestTeser) AssertRules() {
	for _, rule := range r.testrules {
		switch rule.AssertType {
		case ExpectTestHeaderEqual:
			r.AssertHeaderEqual(HeaderType(rule.PropOrExpr), rule.ExpectedVal)
		case ExpectTestHeaderContains:
			r.AssertHeaderContains(HeaderType(rule.PropOrExpr), rule.ExpectedVal)
		case ExpectTestStatusBetween:
			start, end := rule.ExpectedBetweenInt[0], rule.ExpectedBetweenInt[1]
			r.AssertStatusBetween(start, end)
		case ExpectTestBodyContains:
			r.AssertBodyContains(rule.ExpectedVal)
		case ExpectTestBodyJqEqual:
			r.AssertBodyJqEqual(rule.PropOrExpr, rule.ExpectedVal)
		}
	}
}
