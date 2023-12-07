package got

import (
	"bytes"
	"context"
	"io"
	"mime"
	"net/http"
	"path/filepath"
)

// ReqMIME option type, it should be like ".json", "test.json", "a/b/c.jpg", etc
type ReqMIME string

// Req the url. The method is the http method, default value is "GET".
// If an option is http.Header, it will be used as the request header.
// If an option is [Utils.ReqMIME], it will be used to set the Content-Type header.
// If an option is [context.Context], it will be used as the request context.
// Other option type will be treat as request body, it will be encoded by Utils.Write .
func (ut Utils) Req(method, url string, options ...interface{}) *ResHelper {
	ut.Helper()

	header := http.Header{}
	var host string
	var contentType string
	var body io.Reader
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
		default:
			buf := bytes.NewBuffer(nil)
			ut.Write(val)(buf)
			body = buf
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return &ResHelper{ut, nil, err}
	}

	if header != nil {
		req.Header = header
	}

	req.Host = host
	req.Header.Set("Content-Type", contentType)

	res, err := http.DefaultClient.Do(req)
	return &ResHelper{ut, res, err}
}

// ResHelper of the request
type ResHelper struct {
	ut Utils
	*http.Response
	err error
}

// Bytes body
func (res *ResHelper) Bytes() *bytes.Buffer {
	res.ut.Helper()
	return res.ut.Read(res.Body)
}

// String body
func (res *ResHelper) String() string {
	res.ut.Helper()
	return res.Bytes().String()
}

// JSON body
func (res *ResHelper) JSON() (v interface{}) {
	res.ut.Helper()
	return res.ut.JSON(res.Body)
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
