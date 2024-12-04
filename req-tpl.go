package goitest

import "encoding/json"
type RequestTpl struct{
	QueryParam  map[string]string
	FormData    map[string]any
	Json any
	// files       map[string]string     // field -> path
	// fileHeaders map[string]fileHeader // field -> contents
}

func newRequestTpl() *RequestTpl {
	return &RequestTpl{
		QueryParam: make(map[string]string),
		FormData:   make(map[string]any),
	}
}

func (r *RequestTeser) SetQueryParamsTpl(params map[string]string) *RequestTeser {
	for p, v := range params {
		r.SetQueryParamTpl(p, v)
	}
	return r
}
func (r *RequestTeser) SetQueryParamTpl(k, jqExpr string) *RequestTeser {
	r.t.Helper()
	r.tpl.QueryParam[k] = jqExpr
	// segs := strings.SplitN(jqExpr, ".", 2)
		// name, jqExpr2:=segs[0], segs[1]
		vi, err:=integration.GetValue(jqExpr)
		if err != nil {
			r.t.Fatal("Invalid query param tpl: ", jqExpr)
		}
		v, ok := vi.(string)
		if !ok {
			jsonBytes, err := json.Marshal(vi)
			if err != nil {
				r.t.Fatal("Invalid query param tpl: ", jqExpr)
			}
			v = string(jsonBytes)
		}
		r.SetQueryParam(k, v)
		// r.t.Fatal("Invalid query param tpl: ", jqExpr)
	return r
}