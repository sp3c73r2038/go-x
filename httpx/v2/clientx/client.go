package clientx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	urllib "net/url"
	"strings"

	// "github.com/stretchr/testify/assert"
	"github.com/sp3c73r2038/go-x/httpx"
)

type HTTPClientv2 interface {
	Do(string, ...ReqOption) (*Response, error)
	Delete(string, ...ReqOption) (*Response, error)
	Get(string, ...ReqOption) (*Response, error)
	Head(string, ...ReqOption) (*Response, error)
	Option(string, ...ReqOption) (*Response, error)
	Patch(string, ...ReqOption) (*Response, error)
	Post(string, ...ReqOption) (*Response, error)
	Put(string, ...ReqOption) (*Response, error)
	Trace(string, ...ReqOption) (*Response, error)
	Connect(string, ...ReqOption) (*Response, error)

	Cookies(string) ([]*http.Cookie, error)
	SetHeader(k, v string)
	GetClient() *http.Client
}

type SimpleHTTPClient struct {
	auth   httpx.Auth
	client *http.Client
	header http.Header
}

func (this *SimpleHTTPClient) createReq(url string, option *ReqOptions) (rv *http.Request, err error) {

	var body io.Reader

	// prepare the body
	switch option.Method {
	case http.MethodGet:
	case http.MethodHead:
	case http.MethodOptions:
	default: // above are no body
		if option.Body != nil {
			// use provided body first
			body = option.Body
		} else {
			if option.Object != nil {
				var ct = option.Header.Get("Content-Type")
				switch ct {
				case MIME_JSON:
					var buf bytes.Buffer
					var enc = json.NewEncoder(&buf)
					err = enc.Encode(option.Object)
					if err != nil {
						fmt.Errorf("encode json: %w", err)
						return
					}

					body = &buf
				// case "application/x-yaml":
				// case "text/yaml":
				case MIME_FORM:
					var params = option.Object.(map[string]string)
					var form = urllib.Values{}
					for k := range params {
						form.Set(k, params[k])
					}
					body = strings.NewReader(form.Encode())
				default:
				}
			}
		}
	}

	rv, err = http.NewRequest(option.Method, url, body)
	if err != nil {
		fmt.Errorf("create req: %w", err)
		return
	}

	rv.Header = option.Header

	// add any additional arguments
	if len(option.Args) > 0 {
		var qs = rv.URL.Query()
		for k := range option.Args {
			qs.Set(k, option.Args[k])
		}
		rv.URL.RawQuery = qs.Encode()
	}

	return
}

func (this *SimpleHTTPClient) tweak(option *ReqOptions) {
	// inject header from client level
	for k := range this.header {
		WithHeader(k, this.header.Get(k))(option)
	}

	// inject auth header
	if this.auth != nil {
		WithHeader("Authorization", this.auth.AuthHeader())(option)
	}
}

func (this *SimpleHTTPClient) Do(
	url string, opts ...ReqOption) (rv *Response, err error) {

	var option = DefaultReqOptions()

	for _, opt := range opts {
		opt(option)
	}

	this.tweak(option)

	var resp *http.Response
	var req *http.Request

	req, err = this.createReq(url, option)
	if err != nil {
		fmt.Errorf("create req: %w", err)
		return
	}

	resp, err = this.client.Do(req)
	if err != nil {
		fmt.Errorf("do req: %w", err)
		return
	}

	if _, ok := option.Expected[resp.StatusCode]; !ok {
		err = fmt.Errorf("http status %d", resp.StatusCode)
		return
	}

	rv = &Response{
		resp:       resp,
		StatusCode: resp.StatusCode,
	}

	return
}

func (this *SimpleHTTPClient) Get(
	url string, opts ...ReqOption) (rv *Response, err error) {
	return this.Do(url, opts...)
}

func (this *SimpleHTTPClient) Head(
	url string, opts ...ReqOption) (rv *Response, err error) {
	opts = append(opts, WithMethod(http.MethodHead))
	return this.Do(url, opts...)
}

func (this *SimpleHTTPClient) Connect(
	url string, opts ...ReqOption) (rv *Response, err error) {
	opts = append(opts, WithMethod(http.MethodConnect))
	return this.Do(url, opts...)
}

func (this *SimpleHTTPClient) Patch(
	url string, opts ...ReqOption) (rv *Response, err error) {
	opts = append(opts, WithMethod(http.MethodPatch))
	return this.Do(url, opts...)
}

func (this *SimpleHTTPClient) Put(
	url string, opts ...ReqOption) (rv *Response, err error) {
	opts = append(opts, WithMethod(http.MethodPut))
	return this.Do(url, opts...)
}

func (this *SimpleHTTPClient) Post(
	url string, opts ...ReqOption) (rv *Response, err error) {
	opts = append(opts, WithMethod(http.MethodPost))
	return this.Do(url, opts...)
}

func (this *SimpleHTTPClient) Delete(
	url string, opts ...ReqOption) (rv *Response, err error) {
	opts = append(opts, WithMethod(http.MethodDelete))
	return this.Do(url, opts...)
}

func (this *SimpleHTTPClient) Cookies(uri string) (
	rv []*http.Cookie, err error) {
	var u *urllib.URL
	u, err = urllib.Parse(uri)
	if err != nil {
		return
	}
	rv = this.client.Jar.Cookies(u)
	return
}

func (this *SimpleHTTPClient) SetHeader(k, v string) {
	this.header.Set(k, v)
}
func (this *SimpleHTTPClient) GetClient() *http.Client {
	return this.client
}
