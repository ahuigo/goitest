package goitest

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"regexp"
	"testing"
)

type ContentType string

const (
	ContentTypeNone       ContentType = ""
	ContentTypeFormEncode ContentType = "application/x-www-form-urlencoded"
	ContentTypeFormData   ContentType = "multipart/form-data"
	ContentTypeJson       ContentType = "application/json"
	ContentTypePlain      ContentType = "text/plain"
)

type fileHeader struct {
	Filename string
	Header   textproto.MIMEHeader
	Size     int64
	content  []byte
	// tmpfile   string
	// tmpoff    int64
	// tmpshared bool
}

type RequestTeser struct {
	t      *testing.T
	name   string
	rawreq *http.Request
	url    string

	queryParam  url.Values
	formData    url.Values
	isMultiPart bool
	json        any
	files       map[string]string     // field -> path
	fileHeaders map[string]fileHeader // field -> contents
	// request 模板
	tpl *RequestTpl
	// Response
	Response *http.Response
	// Httptest Response
	resp     *httptest.ResponseRecorder
	// testrules
	testrules []TestRule
}

const valid_name_pattern = `^[\w\-]+$`
var validNamePattern = regexp.MustCompile(valid_name_pattern)

func R(t *testing.T, name string) *RequestTeser {
	t.Helper()
	if !validNamePattern.MatchString(name) {
		t.Fatalf("Invalid characters in name: %s", name)
	}
	req :=&RequestTeser{
		t: t,
		name: name,
		tpl: newRequestTpl(),
		rawreq: &http.Request{
			Method:     "GET",
			Header:     make(http.Header),
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
		},
		queryParam: make(url.Values),
		formData:   make(map[string][]string),
		// paramsList:  make(map[string][]string),
		files:       make(map[string]string),
		fileHeaders: make(map[string]fileHeader),
	}
	return req
}
func (r *RequestTeser) Sync() {
	getIntegrationTest().AddReqCase(r)
}

func (r *RequestTeser) Test(name string, f func(*RequestTeser)) {
	r2 := R(r.t, name)
	f(r2)
}

/******************host*************************/
func (r *RequestTeser) SetHost(host string) *RequestTeser {
	r.rawreq.Host = host
	return r
}

/******************header *************************/
func (r *RequestTeser) SetHeader(key, value string) *RequestTeser {
	r.rawreq.Header.Set(key, value)
	return r
}

func (r *RequestTeser) SetAuthBasic(username, password string) *RequestTeser {
	r.rawreq.SetBasicAuth(username, password)
	return r
}

func (r *RequestTeser) SetAuthBearer(token string) *RequestTeser {
	r.rawreq.Header.Set("Authorization", "Bearer "+token)
	return r
}

func (r *RequestTeser) SetContentType(ct ContentType) *RequestTeser {
	r.rawreq.Header.Set("Content-Type", string(ct))
	return r
}

func (r *RequestTeser) AddCookies(cookies []*http.Cookie) *RequestTeser {
	for _, cookie := range cookies {
		r.rawreq.AddCookie(cookie)
	}
	return r
}
func (r *RequestTeser) AddCookieKV(name, value string) *RequestTeser {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
	}
	r.rawreq.AddCookie(cookie)
	return r
}

/************** params **********************/
/************** file **********************/
func (r *RequestTeser) AddFile(fieldname, path string) *RequestTeser {
	r.files[fieldname] = path
	return r
}

func (r *RequestTeser) AddFileHeader(fieldname, filename string, content []byte) *RequestTeser {
	r.fileHeaders[fieldname] = fileHeader{
		Filename: filename,
		content:  content,
		Size:     int64(len(content)),
	}
	return r
}

func (r *RequestTeser) SetUrl(url string) *RequestTeser {
	r.url = url
	return r
}

func (r *RequestTeser) SetReq(method string, url string) *RequestTeser {
	r.rawreq.Method = method
	r.url = url
	return r
}

/************** params **********************/
func (r *RequestTeser) SetQueryParams(params map[string]string) *RequestTeser {
	for p, v := range params {
		r.SetQueryParam(p, v)
	}
	return r
}
func (r *RequestTeser) SetQueryParam(param, value string) *RequestTeser {
	r.queryParam.Set(param, value)
	return r
}

func (r *RequestTeser) SetQueryParamsFromValues(params url.Values) *RequestTeser {
	for p, v := range params {
		for _, pv := range v {
			r.queryParam.Add(p, pv)
		}
	}
	return r
}

/************** body(bytes) **********************/
func (r *RequestTeser) SetBody(body []byte) *RequestTeser {
	r.rawreq.Body = io.NopCloser(bytes.NewReader(body))
	return r
}

/************** body(form) **********************/
// Set Form data(encode or multipart)
func (r *RequestTeser) SetIsMultiPart(b bool) *RequestTeser {
	r.isMultiPart = b
	return r
}
func (r *RequestTeser) SetFormData(data map[string]string) *RequestTeser {
	for k, v := range data {
		r.formData.Set(k, v)
	}
	return r
}

// SetFormDataFromValues method appends multiple form parameters with multi-value
//
//	SetFormDataFromValues(url.Values{"words": []string{"book", "glass", "pencil"},})
func (r *RequestTeser) SetFormDataFromValues(data url.Values) *RequestTeser {
	for k, v := range data {
		for _, kv := range v {
			r.formData.Add(k, kv)
		}
	}
	return r
}

/************** body(json) **********************/
func (r *RequestTeser) SetJson(data any) *RequestTeser {
	r.json = data
	return r
}

/************** body(plain) SetBody(bytes) **********************/

/************** utils **********************/
func (r *RequestTeser) GetRawreq() *http.Request {
	return r.rawreq
}

func (r *RequestTeser) SetCtx(ctx context.Context) *RequestTeser {
	r.rawreq = r.rawreq.WithContext(ctx)
	return r
}

func (r *RequestTeser) EnableTrace(ctx context.Context) *RequestTeser {
	trace := clientTraceNew(r.rawreq.Context())
	r.rawreq = r.rawreq.WithContext(trace.ctx)
	return r
}
