package got

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"path/filepath"
)

// ReqMIME option type, it should be like ".json", "test.json", "a/b/c.jpg", etc
type ReqMIME string

type ReqClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientDoFunc func(req *http.Request) (*http.Response, error)

func (f ClientDoFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Req is a helper method to send http request. It will handle errors automatically, so you don't need to check errors.
// The method is the http method, default value is "GET".
// If an option is [http.Header], it will be used as the request header.
// If an option is [ReqMIME], it will be used to set the Content-Type header.
// If an option is [context.Context], it will be used as the request context.
// If an option is [ReqClient], it will be used as the http client to send the request.
// Other option type will be treat as request body, it will be encoded by [Utils.Write].
// Some request examples:
//
//	Req("GET", "http://example.com")
//	Req("GET", "http://example.com", context.TODO())
//	Req("POST", "http://example.com", map[string]any{"a": 1})
//	Req("POST", "http://example.com", http.Header{"Host": "example.com"}, ReqMIME(".json"), map[string]any{"a": 1})
func (ut Utils) Req(method, url string, options ...interface{}) *ResHelper {
	ut.Helper()

	header := http.Header{}
	var host string
	var contentType string
	var body io.Reader
	var client ReqClient = http.DefaultClient
	ctx := context.Background()

	for _, item := range options {
		switch val := item.(type) {
		case http.Header:
			host = val.Get("Host")
			val.Del("Host")
			header = val
		case ReqMIME:
			contentType = mime.TypeByExtension(filepath.Ext(string(val)))
		case context.Context:
			ctx = val
		case ReqClient:
			client = val
		default:
			buf := bytes.NewBuffer(nil)
			ut.Write(val)(buf)
			body = buf
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return &ResHelper{ut, nil, err, nil}
	}

	if header != nil {
		req.Header = header
	}

	req.Host = host
	req.Header.Set("Content-Type", contentType)

	res, err := client.Do(req)
	return &ResHelper{ut, res, err, nil}
}

// ResHelper of the request
type ResHelper struct {
	ut Utils
	*http.Response
	err error

	read *bytes.Buffer
}

// Bytes parses body as [*bytes.Buffer] and returns the result
func (res *ResHelper) Bytes() *bytes.Buffer {
	res.ut.Helper()
	res.ut.err(res.err)

	if res.read == nil {
		res.read = res.ut.Read(res.Body)
	}

	return res.read
}

// String parses body as string and returns the result
func (res *ResHelper) String() string {
	res.ut.Helper()
	return res.Bytes().String()
}

// JSON parses body as json and returns the result
func (res *ResHelper) JSON() (v interface{}) {
	res.ut.Helper()
	res.ut.err(res.err)
	return res.ut.JSON(res.String())
}

// Unmarshal body to v as json, it's like [json.Unmarshal].
func (res *ResHelper) Unmarshal(v interface{}) {
	res.ut.Helper()
	res.ut.err(json.Unmarshal(res.Bytes().Bytes(), v))
}

// Err of request protocol
func (res *ResHelper) Err() error {
	return res.err
}

func (ut Utils) err(err error) {
	ut.Helper()

	if err != nil {
		ut.Fatal(err)
	}
}

// there no way to stop a blocking test from outside
var panicWithTrace = func(v interface{}) {
	panic(v)
}
