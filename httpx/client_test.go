package httpx

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/sp3c73r2038/go-x/common"
)

func inspectCookie(client *SimpleHTTPClient, uri string) {
	var err error
	common.Logger.Info("cookies: ")
	var cookies []*http.Cookie
	cookies, err = client.Cookies(uri)
	common.Must(err)
	for _, cookie := range cookies {
		common.Logger.Infof("%s: %s", cookie.Name, cookie.Value)
	}
}

func doReq(client *SimpleHTTPClient, uri string) (err error) {
	var resp *http.Response
	resp, err = client.Get(uri, nil, nil)
	common.Must(err)
	defer resp.Body.Close()
	common.Logger.Info(resp.StatusCode)
	var body []byte
	body, err = io.ReadAll(resp.Body)
	common.Must(err)
	common.Logger.Info("body len:", len(string(body)))

	inspectCookie(client, uri)

	return
}

func TestHTTPClient(t *testing.T) {
	var err error
	var u *url.URL

	u, err = url.Parse("http://127.0.0.1:3128")
	common.Must(err)

	var auth = NewBasicAuth("", "")

	var client = NewHTTPClientBuilder().SetAuth(auth).SetProxy(u).Build()
	//
	doReq(client, "http://127.0.0.1")
	//
	doReq(client, "https://www.google.com")
	//
	doReq(client, "https://www.google.com")
	//
	client = NewHTTPClientBuilder().SetSkipTLSVerify(true).Build()
	// doReq(client, "https://10.10.0.1:5443")
	var url = "https://10.10.0.1:5443/logincheck"
	var resp *http.Response
	resp, err = client.PostForm(url, map[string]string{
		"username":  "lei",
		"secretkey": "",
	}, nil, nil)
	common.Must(err)
	defer resp.Body.Close()
	common.Logger.Info(resp.StatusCode)
	var body []byte
	body, err = io.ReadAll(resp.Body)
	common.Must(err)
	common.Logger.Info(string(body))

}
