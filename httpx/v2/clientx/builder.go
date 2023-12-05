package clientx

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"net/http"
	"net/http/cookiejar"
	urllib "net/url"
	"os"
	"time"

	// "github.com/stretchr/testify/assert"
	"github.com/sp3c73r2038/go-x/common"
	"github.com/sp3c73r2038/go-x/httpx"
)

type HTTPClientBuilder struct {
	timeout time.Duration
	proxy   *urllib.URL
	skipTLS bool
	auth    httpx.Auth
	rootCAs *x509.CertPool
}

func NewHTTPClientBuilder() *HTTPClientBuilder {
	return &HTTPClientBuilder{
		timeout: time.Second * 60,
	}
}

func NewCertPoolFromPEM(b []byte) (rv *x509.CertPool, err error) {
	var ok bool

	rv = x509.NewCertPool()
	ok = rv.AppendCertsFromPEM(b)
	if !ok {
		err = errors.New("can't append certs from pem")
	}
	return
}

func NewCertPoolFromFile(in string) (rv *x509.CertPool, err error) {
	var ok = false

	rv = x509.NewCertPool()

	var b []byte
	b, err = os.ReadFile(in)
	if err != nil {
		return
	}

	ok = rv.AppendCertsFromPEM(b)
	if !ok {
		err = errors.New("can't append certs from pem")
	}
	return
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

func (this *HTTPClientBuilder) SetCA(ca *x509.CertPool) (rv *HTTPClientBuilder) {
	this.rootCAs = ca
	return this
}

func (this *HTTPClientBuilder) Build() (rv *SimpleHTTPClient) {
	var err error
	var jar *cookiejar.Jar
	jar, err = cookiejar.New(nil)

	common.Must(err)

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
					RootCAs:            this.rootCAs,
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
