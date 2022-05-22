package clientx

import (
	"encoding/json"
	"io"
	"net/http"
	// "github.com/stretchr/testify/assert"
)

type Response struct {
	resp       *http.Response
	StatusCode int
}

func (this *Response) GetRaw() *http.Response {
	return this.resp
}

func (this *Response) JSON(rt interface{}) (err error) {
	var dec = json.NewDecoder(this.resp.Body)
	err = dec.Decode(&rt)
	return
}

func (this *Response) Bytes() (rv []byte, err error) {
	rv, err = io.ReadAll(this.resp.Body)
	return
}

func (this *Response) Text() (rv string, err error) {
	var b []byte
	b, err = io.ReadAll(this.resp.Body)
	rv = string(b)
	return
}

func (this *Response) Close() {
	if this.resp != nil {
		this.resp.Body.Close()
	}
}
