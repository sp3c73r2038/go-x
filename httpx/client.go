package httpx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"net/http/cookiejar"
	urllib "net/url"
	"strings"
	"time"

	"github.com/sp3c73r2038/go-x/common"
)

type (
	HTTPClient interface {
		Get(url string, args map[string]string, headers map[string]string) (
			*http.Response, error)

		Post(url string, body []byte, args map[string]string,
			headers map[string]string) (*http.Response, error)

		PostForm(url string, body map[string]string, args map[string]string,
			headers map[string]string) (*http.Response, error)

		PostJSON(url string, body interface{}, args map[string]string,
			headers map[string]string) (*http.Response, error)

		PutJSON(url string, body interface{}, args map[string]string,
			headers map[string]string) (*http.Response, error)

		Delete(url string, args map[string]string, headers map[string]string) (
			*http.Response, error)

		Req(method string, url string, body []byte, args map[string]string,
			headers map[string]string) (*http.Response, error)

		Cookies(uri string) ([]*http.Cookie, error)

		SetHeader(k, v string)
	}

	SimpleHTTPClient struct {
		httpclient *http.Client
		auth       Auth
		header     http.Header
	}
)

func (this *SimpleHTTPClient) GetRawClient() (rv *http.Client) {
	return this.httpclient
}

func (this *SimpleHTTPClient) Cookies(uri string) (rv []*http.Cookie, err error) {
	var u *urllib.URL
	u, err = urllib.Parse(uri)
	if err != nil {
		return
	}
	rv = this.httpclient.Jar.Cookies(u)
	return
}

func (this *SimpleHTTPClient) SetHeader(k, v string) {
	this.header.Set(k, v)
	return
}

func (this *SimpleHTTPClient) Get(
	url string, args map[string]string, headers map[string]string) (
	rv *http.Response, err error) {

	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	if args != nil {
		var q = req.URL.Query()
		for k, v := range args {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return this.sendReq(req)
}

func (this *SimpleHTTPClient) Post(url string, body []byte, args map[string]string,
	headers map[string]string) (rv *http.Response, err error) {

	var req *http.Request
	req, err = http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return
	}

	if args != nil {
		var q = req.URL.Query()
		for k, v := range args {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return this.sendReq(req)
}

func (this *SimpleHTTPClient) PostForm(url string, body map[string]string, args map[string]string,
	headers map[string]string) (rv *http.Response, err error) {

	var form = urllib.Values{}
	for k, v := range body {
		form.Set(k, v)
	}

	var req *http.Request

	req, err = http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return
	}

	if args != nil {
		var q = req.URL.Query()
		for k, v := range args {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return this.sendReq(req)
}

func (this *SimpleHTTPClient) PostJSON(url string, v interface{}, args map[string]string,
	headers map[string]string) (rv *http.Response, err error) {

	var buf bytes.Buffer
	var req *http.Request

	err = json.NewEncoder(&buf).Encode(&v)
	if err != nil {
		return
	}

	req, err = http.NewRequest("POST", url, &buf)
	if err != nil {
		return
	}

	if args != nil {
		var q = req.URL.Query()
		for k, v := range args {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")

	return this.sendReq(req)
}

func (this *SimpleHTTPClient) PutJSON(url string, v interface{}, args map[string]string,
	headers map[string]string) (rv *http.Response, err error) {

	var buf bytes.Buffer
	var req *http.Request

	err = json.NewEncoder(&buf).Encode(&v)
	if err != nil {
		return
	}

	req, err = http.NewRequest("PUT", url, &buf)
	if err != nil {
		return
	}

	if args != nil {
		var q = req.URL.Query()
		for k, v := range args {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")

	return this.sendReq(req)
}

func (this *SimpleHTTPClient) Delete(
	url string, args map[string]string, headers map[string]string) (
	rv *http.Response, err error) {

	var req *http.Request
	req, err = http.NewRequest("DELETE", url, nil)
	if err != nil {
		return
	}

	if args != nil {
		var q = req.URL.Query()
		for k, v := range args {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return this.sendReq(req)
}

func (this *SimpleHTTPClient) Req(method string, url string, body []byte, args map[string]string,
	headers map[string]string) (rv *http.Response, err error) {
	var req *http.Request
	req, err = http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return
	}

	if args != nil {
		var q = req.URL.Query()
		for k, v := range args {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	this.sendReq(req)
	return
}

func (this *SimpleHTTPClient) sendReq(req *http.Request) (rv *http.Response, err error) {

	if len(this.header) > 0 {
		// TODO: when to use Header.Add?
		for k := range this.header {
			req.Header.Set(k, this.header.Get(k))
		}
	}

	if this.auth != nil {
		req.Header.Set("Authorization", this.auth.AuthHeader())
	}

	common.Logger.Debug("req header: ", common.Pretty(req.Header))

	rv, err = this.httpclient.Do(req)
	return
}

func NewHTTPClient() *http.Client {
	return &http.Client{
		// Timeout: timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				KeepAlive: 0,
			}).DialContext,
		},
		CheckRedirect: func(
			req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func NewSimpleHTTPClient() *SimpleHTTPClient {
	return &SimpleHTTPClient{
		httpclient: NewHTTPClient(),
	}
}

type HTTPClientBuilder struct {
	timeout time.Duration
	proxy   *urllib.URL
	skipTLS bool
	auth    Auth
}

func NewHTTPClientBuilder() *HTTPClientBuilder {
	return &HTTPClientBuilder{
		timeout: time.Second * 60,
		skipTLS: false,
	}
}

func (this *HTTPClientBuilder) SetTimeout(
	timeout time.Duration) (rv *HTTPClientBuilder) {
	this.timeout = timeout
	return this
}

func (this *HTTPClientBuilder) SetAuth(
	auth Auth) (rv *HTTPClientBuilder) {
	this.auth = auth
	return this
}

func (this *HTTPClientBuilder) SetProxy(
	proxy *urllib.URL) (rv *HTTPClientBuilder) {
	this.proxy = proxy
	return this
}

func (this *HTTPClientBuilder) SetSkipTLSVerify(
	skip bool) (rv *HTTPClientBuilder) {
	this.skipTLS = skip
	return this
}

func (this *HTTPClientBuilder) Build() (rv *SimpleHTTPClient) {
	var err error
	var jar *cookiejar.Jar
	jar, err = cookiejar.New(nil)

	common.Must(err)

	rv = &SimpleHTTPClient{
		auth: this.auth,
		httpclient: &http.Client{
			Jar:     jar,
			Timeout: this.timeout,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					KeepAlive: 0,
				}).DialContext,
				Proxy: http.ProxyURL(this.proxy),
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: this.skipTLS,
				},
			},
			CheckRedirect: func(
				req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		header: http.Header{},
	}
	return rv
}
