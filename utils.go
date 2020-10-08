package got

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// Context helper
type Context struct {
	context.Context
	Cancel func()
}

// Utils for commonly used methods
type Utils struct {
	Testable
}

// Fatal is the same as testing.common.Fatal
func (ut Utils) Fatal(args ...interface{}) {
	ut.Helper()
	ut.Log(args...)
	ut.FailNow()
}

// Fatalf is the same as testing.common.Fatalf
func (ut Utils) Fatalf(format string, args ...interface{}) {
	ut.Helper()
	ut.Logf(format, args...)
	ut.FailNow()
}

// Log is the same as testing.common.Log
func (ut Utils) Log(args ...interface{}) {
	ut.Helper()
	ut.Logf("%s", fmt.Sprintln(args...))
}

// Error is the same as testing.common.Error
func (ut Utils) Error(args ...interface{}) {
	ut.Helper()
	ut.Log(args...)
	ut.Fail()
}

// Errorf is the same as testing.common.Errorf
func (ut Utils) Errorf(format string, args ...interface{}) {
	ut.Helper()
	ut.Logf(format, args...)
	ut.Fail()
}

// Skipf is the same as testing.common.Skipf
func (ut Utils) Skipf(format string, args ...interface{}) {
	ut.Helper()
	ut.Logf(format, args...)
	ut.SkipNow()
}

// Skip is the same as testing.common.Skip
func (ut Utils) Skip(args ...interface{}) {
	ut.Helper()
	ut.Log(args...)
	ut.SkipNow()
}

// Parallel is the same as testing.T.Parallel
func (ut Utils) Parallel() {
	reflect.ValueOf(ut.Testable).MethodByName("Parallel").Call(nil)
}

// Context that will be canceled after the test
func (ut Utils) Context() Context {
	ctx, cancel := context.WithCancel(context.Background())
	ut.Cleanup(cancel)
	return Context{ctx, cancel}
}

// Timeout context that will be canceled after the test
func (ut Utils) Timeout(d time.Duration) Context {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	ut.Cleanup(cancel)
	return Context{ctx, cancel}
}

// Srand generates a random string with the specified length
func (ut Utils) Srand(l int) string {
	ut.Helper()
	b := make([]byte, (l+1)/2)
	_, err := rand.Read(b)
	ut.err(err)
	return hex.EncodeToString(b)[:l]
}

// Open a file. Override it if create is true. Directories will be auto-created.
// path will be joined with filepath.Join so that it's cross-platform
func (ut Utils) Open(create bool, path ...string) (f *os.File) {
	p := filepath.Join(path...)

	dir := filepath.Dir(p)
	_ = os.MkdirAll(dir, 0755)

	var err error
	if create {
		f, err = os.Create(p)
	} else {
		f, err = os.Open(p)
	}
	ut.err(err)
	return f
}

// Read all from r
func (ut Utils) Read(r io.Reader) []byte {
	ut.Helper()
	b, err := ioutil.ReadAll(r)
	ut.err(err)
	return b
}

// ReadString from r
func (ut Utils) ReadString(r io.Reader) string {
	ut.Helper()
	return string(ut.Read(r))
}

// JSON from string, []byte, or io.Reader
func (ut Utils) JSON(src interface{}) (v interface{}) {
	ut.Helper()

	var b []byte
	switch obj := src.(type) {
	case []byte:
		b = obj
	case string:
		b = []byte(obj)
	case io.Reader:
		var err error
		b, err = ioutil.ReadAll(obj)
		ut.err(err)
	}
	ut.err(json.Unmarshal(b, &v))
	return
}

// Write obj to the writer. Encode obj to []byte and cache it for writer.
// If obj is not []byte, string, or io.Reader, it will be encoded as JSON.
func (ut Utils) Write(obj interface{}) (writer func(io.Writer)) {
	var cache io.ReadWriter
	return func(w io.Writer) {
		ut.Helper()

		if cache != nil {
			_, err := io.Copy(w, cache)
			ut.err(err)
			return
		}

		cache = bytes.NewBuffer(nil)
		w = io.MultiWriter(cache, w)

		var err error
		switch v := obj.(type) {
		case []byte:
			_, err = w.Write(v)
		case string:
			_, err = w.Write([]byte(v))
		case io.Reader:
			_, err = io.Copy(w, v)
		default:
			err = json.NewEncoder(w).Encode(v)
		}
		ut.err(err)
	}
}

// HandleHTTP handles a request. If file exists serve the file content. The file will be used to set the Content-Type header.
// If the file doesn't exist, the value will be encoded by G.Write(value) and used as the response body.
func (ut Utils) HandleHTTP(file string, value ...interface{}) func(http.ResponseWriter, *http.Request) {
	var obj interface{}
	if len(value) > 1 {
		obj = value
	} else if len(value) == 1 {
		obj = value[0]
	}

	write := ut.Write(obj)

	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(file); err == nil {
			http.ServeFile(w, r, file)
			return
		}

		if obj == nil {
			return
		}

		w.Header().Add("Content-Type", mime.TypeByExtension(filepath.Ext(file)))

		write(w)
	}
}

// Serve http on a random port. The server will be auto-closed after the test.
func (ut Utils) Serve() *Router {
	ut.Helper()

	mux := http.NewServeMux()
	srv := &http.Server{Handler: mux}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	ut.err(err)

	ut.Cleanup(func() { ut.err(srv.Close()) })

	go func() { _ = srv.Serve(l) }()

	u, err := url.Parse("http://" + l.Addr().String())
	ut.err(err)

	return &Router{ut, u, srv, mux}
}

// Router of a http server
type Router struct {
	ut      Utils
	HostURL *url.URL
	Server  *http.Server
	Mux     *http.ServeMux
}

// URL will prefix the path with the server's host
func (rt *Router) URL(path ...string) string {
	return rt.HostURL.String() + strings.Join(path, "")
}

// Route on the pattern. Check the doc of http.ServeMux for the syntax of pattern.
// It will use G.HandleHTTP to handle each request.
func (rt *Router) Route(pattern, file string, value ...interface{}) *Router {
	h := rt.ut.HandleHTTP(file, value...)

	rt.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})

	return rt
}

// Req the url. The method is the http method. The body will be encoded by G.Write(body) .
// When the len(body) is greater than 2, the first item should be a file extension string for the Content-Type header,
// such as ".json", ".jpg".
func (ut Utils) Req(method, url string, body ...interface{}) *ResHelper {
	ut.Helper()

	var r io.Reader

	var obj interface{}
	file := ""
	switch len(body) {
	case 0:
	case 1:
		obj = body[0]
	case 2:
		file = body[0].(string)
		obj = body[1]
	default:
		file = body[0].(string)
		obj = body[1:]
	}

	if obj != nil {
		var w io.WriteCloser
		r, w = io.Pipe()
		go func() {
			ut.Write(obj)(w)
			ut.err(w.Close())
		}()
	}

	req, err := http.NewRequest(method, url, r)
	ut.err(err)

	req.Header.Add("Content-Type", mime.TypeByExtension(filepath.Ext(file)))

	res, err := http.DefaultClient.Do(req)
	ut.err(err)

	return &ResHelper{ut, res}
}

// ResHelper of the request
type ResHelper struct {
	ut Utils
	*http.Response
}

// Bytes body
func (res *ResHelper) Bytes() []byte {
	res.ut.Helper()
	return res.ut.Read(res.Body)
}

// String body
func (res *ResHelper) String() string {
	res.ut.Helper()
	return string(res.Bytes())
}

// JSON body
func (res *ResHelper) JSON() (v interface{}) {
	res.ut.Helper()
	res.ut.err(json.Unmarshal(res.Bytes(), &v))
	return
}

func (ut Utils) err(err error) {
	if err != nil {
		ut.Fatal(err)
	}
}
