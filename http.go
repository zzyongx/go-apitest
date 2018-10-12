package apitest

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"runtime/debug"
	"strings"
	"testing"
	"time"
)

type Api struct {
	Endpoint string
	Timeout  time.Duration
	Headers  http.Header
	Cookies  map[string]string
	Request  *ApiRequest
}

type ApiRequest struct {
	Timeout    time.Duration
	Method     string
	Path       string
	PathParams []string
	Params     url.Values

	Headers http.Header
	Cookies map[string]string

	Forms url.Values
	Json  string
}

func NewApiRequest() *ApiRequest {
	return &ApiRequest{
		Params:  url.Values{},
		Headers: http.Header{},
		Cookies: make(map[string]string, 0),
		Forms:   url.Values{},
	}
}

func NewHttpTest(endpoint string) *Api {
	return &Api{
		Endpoint: endpoint,
		Timeout:  time.Second,
		Headers:  http.Header{},
		Cookies:  make(map[string]string, 0),
		Request:  NewApiRequest(),
	}
}

func (this *Api) Get(path string, params ...interface{}) *Api {
	this.Request.Method = "GET"
	return this.request(path, params...)
}

func (this *Api) Param(key string, value interface{}) *Api {
	this.Request.Params.Add(key, MustInterfaceToString(value))
	return this
}

func (this *Api) Post(path string, params ...interface{}) *Api {
	this.Request.Method = "POST"
	return this.request(path, params...)
}

func (this *Api) Put(path string, params ...interface{}) *Api {
	this.Request.Method = "PUT"
	return this.request(path, params...)
}

func (this *Api) Form(key string, value interface{}) *Api {
	this.Request.Forms.Add(key, MustInterfaceToString(value))
	return this
}

func (this *Api) Delete(path string, params ...interface{}) *Api {
	this.Request.Method = "DELETE"
	return this.request(path, params...)
}

func (this *Api) request(path string, params ...interface{}) *Api {
	this.Request.Path = path
	for _, param := range params {
		this.Request.PathParams = append(this.Request.PathParams, MustInterfaceToString(param))
	}
	return this
}

func (this *Api) Header(name, value string) *Api {
	this.Request.Headers.Add(name, value)
	return this
}

func (this *Api) Cookie(name, value string) *Api {
	this.Request.Cookies[name] = value
	return this
}

func (this *Api) Json(json string) *Api {
	this.Request.Json = json
	return this
}

type ApiExpect struct {
	t             *testing.T
	req           string
	statusExpect  *StatusExpect
	headersExpect *HeadersExpect
	cookiesExpect *CookiesExpect
	jsonExpect    *JsonExpect
	plainExpect   *PlainExpect
}

func (this *Api) consUrl() string {
	path := this.Request.Path
	pathParams := this.Request.PathParams
	if len(pathParams) > 0 {
		re := regexp.MustCompile("({[^{]+})")
		path = re.ReplaceAllStringFunc(path, func(x string) string {
			repl := pathParams[0]
			if len(pathParams) > 1 {
				pathParams = pathParams[1:]
			}
			return repl
		})
	}
	if len(this.Request.Params) > 0 {
		params := this.Request.Params.Encode()
		if strings.Index(path, "?") > 0 {
			path = path + "&" + params
		} else {
			path = path + "?" + params
		}
	}
	return this.Endpoint + path
}

func (this *Api) doRequest(t *testing.T) (string, *http.Response) {
	uri := this.consUrl()
	var contentType string
	var body string
	if this.Request.Method == "POST" || this.Request.Method == "PUT" {
		if len(this.Request.Forms) != 0 {
			body = this.Request.Forms.Encode()
			contentType = "application/x-www-form-urlencoded"
		} else if len(this.Request.Json) > 0 {
			body = this.Request.Json
			contentType = "application/json"
		}
	}

	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = strings.NewReader(body)
	}

	req, err := http.NewRequest(this.Request.Method, uri, bodyReader)
	if err != nil {
		t.Fatalf("NewRequest %v error: %s", this.Request, err)
	}

	for key, values := range this.Headers {
		if _, ok := this.Request.Headers[key]; !ok {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	for key, values := range this.Request.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	req.Header.Set("Content-Type", contentType)

	for name, value := range this.Cookies {
		if _, ok := this.Request.Cookies[name]; !ok {
			req.AddCookie(&http.Cookie{Name: name, Value: value})
		}
	}

	for name, value := range this.Request.Cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: false},
		DisableKeepAlives: true,
	}

	timeout := this.Timeout
	if this.Request.Timeout > 0 {
		timeout = this.Request.Timeout
	}

	reqStr := this.resetRequest(req, body)

	client := &http.Client{Transport: tr, Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("req: %s > error: %s", reqStr, err)
	}
	return reqStr, resp
}

func (this *Api) resetRequest(req *http.Request, body string) string {
	this.Request = NewApiRequest()

	return fmt.Sprintf(`
Method: %s
URL: %s
Headers: %v
Body: %s
`, req.Method, req.URL.String(), req.Header, body)
}

func (this *Api) Expect(t *testing.T) *ApiExpect {
	reqStr, resp := this.doRequest(t)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("req: %s > read response body error: %s", reqStr, err)
	}

	apiExpect := ApiExpect{
		t:   t,
		req: reqStr,
	}

	apiExpect.statusExpect = &StatusExpect{
		ApiExpect: &apiExpect,
		status:    resp.StatusCode,
	}

	apiExpect.headersExpect = &HeadersExpect{
		ApiExpect: &apiExpect,
		headers:   resp.Header,
	}

	apiExpect.cookiesExpect = &CookiesExpect{
		ApiExpect: &apiExpect,
		cookies:   resp.Cookies(),
	}

	apiExpect.jsonExpect = &JsonExpect{
		ApiExpect: &apiExpect,
		data:      string(body),
	}

	apiExpect.plainExpect = &PlainExpect{
		ApiExpect: &apiExpect,
		data:      body,
	}

	return &apiExpect
}

func (this *ApiExpect) Fatalf(format string, args ...interface{}) {
	newArgs := make([]interface{}, 0)
	newArgs = append(newArgs, this.req)
	newArgs = append(newArgs, args...)

	top := -1
	stacks := strings.Split(string(debug.Stack()), "\n")
	for i := len(stacks) - 1; i >= 0; i = i - 1 {
		if ok, _ := regexp.MatchString("^\t.+github.com/zzyongx/go-apitest/[^/]+\\.go", stacks[i]); ok {
			top = i + 1
			break
		}
	}
	stacks = stacks[top:]
	fmt.Println(strings.Join(stacks, "\n"))

	this.t.Fatalf("req: %s > "+format, newArgs...)
}

func (this *ApiExpect) Req() string {
	return this.req
}

func (this *ApiExpect) Status() *StatusExpect {
	return this.statusExpect
}

func (this *ApiExpect) Headers() *HeadersExpect {
	return this.headersExpect
}

func (this *ApiExpect) Cookies(name string) *CookiesExpect {
	this.cookiesExpect.name = name
	return this.cookiesExpect
}

func (this *ApiExpect) Json() *JsonExpect {
	var obj interface{}
	if err := json.Unmarshal([]byte(this.jsonExpect.data), &obj); err != nil {
		this.t.Fatalf("req: %s > parse json %s error: %s", this.req, this.jsonExpect.data, err)
	}
	this.jsonExpect.obj = obj
	return this.jsonExpect
}

func (this *ApiExpect) Plain() *PlainExpect {
	return this.plainExpect
}
