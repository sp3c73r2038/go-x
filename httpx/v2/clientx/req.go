package clientx

import (
	"bytes"
	"io"
	"net/http"
)

type ReqOption = func(*ReqOptions)

type ReqOptions struct {
	Method   string
	Args     map[string]string
	Header   http.Header
	Body     io.Reader
	Object   interface{}
	Expected map[int]*struct{}
}

func DefaultReqOptions() *ReqOptions {
	return &ReqOptions{
		Method: http.MethodGet,
		Args:   make(map[string]string),
		Header: http.Header{},
		Expected: map[int]*struct{}{
			200: nil,
			201: nil,
			204: nil,
		},
	}
}

// type RespOption = func(*http.Response)

func WithMethod(method string) func(*ReqOptions) {
	return func(opt *ReqOptions) {
		opt.Method = method
	}
}

func WithBody(b io.Reader) func(*ReqOptions) {
	return func(opt *ReqOptions) {
		opt.Body = b
	}
}

func WithHeaders(headers map[string]string) func(*ReqOptions) {
	return func(opt *ReqOptions) {
		for k := range headers {
			opt.Header.Set(k, headers[k])
		}
	}
}

func WithHeader(k, v string) func(*ReqOptions) {
	return func(opt *ReqOptions) {
		opt.Header.Set(k, v)
	}
}

func WithBodyObject(ct string, obj interface{}) func(*ReqOptions) {

	return func(opt *ReqOptions) {
		opt.Header.Set("Content-Type", ct)
		opt.Object = obj
	}
}

func WithBytes(b []byte) func(*ReqOptions) {
	return func(opt *ReqOptions) {
		opt.Body = bytes.NewBuffer(b)
	}
}

func WithExpect(codes ...int) func(*ReqOptions) {
	return func(opt *ReqOptions) {
		var ex = make(map[int]*struct{})
		for _, code := range codes {
			ex[code] = nil
		}
		opt.Expected = ex
	}
}

func WithJSON(obj interface{}) func(*ReqOptions) {
	return WithBodyObject(MIME_JSON, obj)
}

func WithForm(obj interface{}) func(*ReqOptions) {
	return WithBodyObject(MIME_FORM, obj)
}

func WithArg(k, v string) func(*ReqOptions) {
	return func(opt *ReqOptions) {
		opt.Args[k] = v
	}
}

func WithArgs(args map[string]string) func(*ReqOptions) {
	return func(opt *ReqOptions) {
		for k := range args {
			opt.Args[k] = args[k]
		}
	}
}
