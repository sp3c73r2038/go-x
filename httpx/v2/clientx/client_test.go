package clientx

import (
	"io"
	"net/http"
	urllib "net/url"
	"testing"
	"time"

	"github.com/sp3c73r2038/go-x/httpx"
	"github.com/stretchr/testify/assert"
	// "snippet/common"
)

func TestBuild(t *testing.T) {
	var err error
	var client = NewHTTPClientBuilder().SetTimeout(time.Second * 10).Build()
	var resp *Response
	resp, err = client.Get("http://localhost")
	if err != nil {
		t.Errorf("get error")
		return
	}
	defer resp.Close()
}

func TestCreateReq(t *testing.T) {
	var err error
	var client = NewHTTPClientBuilder().SetTimeout(time.Second * 10).Build()
	var option = DefaultReqOptions()
	var req *http.Request
	var b []byte
	var q urllib.Values

	// 0. empty req
	req, err = client.createReq("", option)
	if err != nil {
		t.Errorf(": %v", err)
		return
	}
	assert.NotNil(t, req)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Nil(t, req.Body)
	assert.Equal(t, 0, len(req.Header))

	// 1. with header, with method
	WithHeader("H", "V")(option)
	WithMethod(http.MethodHead)(option)
	req, err = client.createReq("", option)
	if err != nil {
		t.Errorf(": %v", err)
		return
	}
	assert.NotNil(t, req)
	assert.Equal(t, http.MethodHead, req.Method)
	assert.Nil(t, req.Body)
	assert.Equal(t, 1, len(req.Header))
	assert.Equal(t, "V", req.Header.Get("H"))

	// // 2. with form
	// option = DefaultReqOptions()
	// WithHeader("H", "V")(option)
	// WithMethod(http.MethodPut)(option)
	// WithForm(map[string]string{"p1": "v1"})(option)
	// req, err = client.createReq("", option)
	// if err != nil {
	// 	t.Errorf(": %v", err)
	// 	return
	// }
	// assert.NotNil(t, req)
	// assert.Equal(t, http.MethodPut, req.Method)
	// assert.NotNil(t, req.Body)
	// assert.Equal(t, MIME_FORM, req.Header.Get("Content-Type"))
	// b, err = io.ReadAll(req.Body)
	// if err != nil {
	// 	t.Errorf(": %v", err)
	// 	return
	// }
	// assert.Equal(t, []byte("p1=v1"), b)

	// 2. with json
	option = DefaultReqOptions()
	WithHeader("H", "V")(option)
	WithMethod(http.MethodPost)(option)
	WithJSON(map[string]string{"p1": "v1"})(option)
	req, err = client.createReq("", option)
	if err != nil {
		t.Errorf(": %v", err)
		return
	}
	assert.NotNil(t, req)
	assert.Equal(t, http.MethodPost, req.Method)
	assert.NotNil(t, req.Body)
	assert.Equal(t, MIME_JSON, req.Header.Get("Content-Type"))
	b, err = io.ReadAll(req.Body)
	if err != nil {
		t.Errorf(": %v", err)
		return
	}
	assert.Equal(t, []byte("{\"p1\":\"v1\"}\n"), b)

	// 3. with arguments to existing url
	option = DefaultReqOptions()
	WithArg("c", "3")(option)
	req, err = client.createReq("http://example.com/path?a=1&b=2", option)
	if err != nil {
		t.Errorf(": %v", err)
		return
	}
	assert.NotNil(t, req)
	assert.NotNil(t, req.URL)
	q = req.URL.Query()
	assert.Equal(t, "1", q.Get("a"))
	assert.Equal(t, "2", q.Get("b"))
	assert.Equal(t, "3", q.Get("c"))

	// 4. client level header
	option = DefaultReqOptions()
	client.SetHeader("H", "V")
	client.tweak(option)
	req, err = client.createReq("http://example.com/path?a=1&b=2", option)
	if err != nil {
		t.Errorf(": %v", err)
		return
	}
	assert.Equal(t, "V", req.Header.Get("H"))

	// 5. client level auth
	option = DefaultReqOptions()
	client.auth = httpx.NewBearerAuth("bearertoken")
	client.tweak(option)
	req, err = client.createReq("http://example.com/path?a=1&b=2", option)
	if err != nil {
		t.Errorf(": %v", err)
		return
	}
	assert.Equal(t, "Bearer bearertoken", req.Header.Get("Authorization"))
}
