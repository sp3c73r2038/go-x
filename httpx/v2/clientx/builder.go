package clientx

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	urllib "net/url"
	"time"

	// "github.com/stretchr/testify/assert"
	"github.com/sp3c73r2038/go-x/httpx"
)

type HTTPClientBuilder struct {
	timeout time.Duration
	proxy   *urllib.URL
	skipTLS bool
	auth    httpx.Auth
}

func NewHTTPClientBuilder() *HTTPClientBuilder {
	return &HTTPClientBuilder{
		timeout: time.Second * 60,
	}
}

func (this *HTTPClientBuilder) SetTimeout(
	timeout time.Duration) (rv *HTTPClientBuilder) {
	this.timeout = timeout
	return this
}

func (this *HTTPClientBuilder) SetAuth(
	auth httpx.Auth) (rv *HTTPClientBuilder) {
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

	httpx.Must(err)

	rv = &SimpleHTTPClient{
		auth: this.auth,
		client: &http.Client{
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
