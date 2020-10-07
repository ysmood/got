package got

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Srand generates a random string with the specified length
func (hp G) Srand(l int) string {
	hp.Helper()
	b := make([]byte, (l+1)/2)
	_, err := rand.Read(b)
	hp.E(err)
	return hex.EncodeToString(b)[:l]
}

// Open a file. Override it if create is true. Directories will be auto-created.
// path will be joined with filepath.Join so that it's cross-platform
func (hp G) Open(create bool, path ...string) (f *os.File) {
	p := filepath.Join(path...)

	dir := filepath.Dir(p)
	_ = os.MkdirAll(dir, 0755)

	var err error
	if create {
		f, err = os.Create(p)
	} else {
		f, err = os.Open(p)
	}
	hp.E(err)
	return f
}

// Read all from r
func (hp G) Read(r io.Reader) []byte {
	hp.Helper()
	b, err := ioutil.ReadAll(r)
	hp.E(err)
	return b
}

// ReadString from r
func (hp G) ReadString(r io.Reader) string {
	hp.Helper()
	return string(hp.Read(r))
}

// ReadJSON from r
func (hp G) ReadJSON(r io.Reader) (v interface{}) {
	hp.Helper()
	hp.E(json.Unmarshal(hp.Read(r), &v))
	return
}

// Write obj to the w
func (hp G) Write(obj interface{}) func(io.Writer) {
	var cache io.ReadWriter
	return func(w io.Writer) {
		hp.Helper()

		if cache != nil {
			hp.E(io.Copy(w, cache))
			return
		}

		cache = bytes.NewBuffer(nil)
		w = io.MultiWriter(cache, w)

		var err error
		switch v := obj.(type) {
		case []byte:
			_, err = w.Write(v)
		case io.Reader:
			_, err = io.Copy(w, v)
		default:
			err = json.NewEncoder(w).Encode(v)
		}
		hp.E(err)
	}
}

// HandleHTTP handles a request
func (hp G) HandleHTTP(file string, value ...interface{}) func(http.ResponseWriter, *http.Request) {
	var obj interface{}
	if len(value) > 1 {
		obj = value
	} else if len(value) == 1 {
		obj = value[0]
	}

	write := hp.Write(obj)

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
func (hp G) Serve() *Router {
	hp.Helper()

	mux := http.NewServeMux()
	srv := &http.Server{Handler: mux}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	hp.E(err)

	hp.Cleanup(func() { hp.E(srv.Close()) })

	go func() { _ = srv.Serve(l) }()

	u, err := url.Parse("http://" + l.Addr().String())
	hp.E(err)

	return &Router{hp, u, srv, mux}
}

// Router of a http server
type Router struct {
	hp      G
	HostURL *url.URL
	Server  *http.Server
	Mux     *http.ServeMux
}

// URL will prefix the path with the server's host
func (rt *Router) URL(path ...string) string {
	return rt.HostURL.String() + "/" + strings.TrimLeft(strings.Join(path, "/"), "/")
}

// Route on the pattern.
func (rt *Router) Route(pattern, file string, value ...interface{}) *Router {
	h := rt.hp.HandleHTTP(file, value...)

	rt.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})

	return rt
}

// Req to the url
func (hp G) Req(method, url string, body ...interface{}) *ResHelper {
	hp.Helper()

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
			hp.Write(obj)(w)
			hp.E(w.Close())
		}()
	}

	req, err := http.NewRequest(method, url, r)
	hp.E(err)

	req.Header.Add("Content-Type", mime.TypeByExtension(filepath.Ext(file)))

	res, err := http.DefaultClient.Do(req)
	hp.E(err)

	return &ResHelper{hp, res}
}

// ResHelper of the request
type ResHelper struct {
	hp G
	*http.Response
}

// Bytes body
func (res *ResHelper) Bytes() []byte {
	res.hp.Helper()
	return res.hp.Read(res.Body)
}

// String body
func (res *ResHelper) String() string {
	res.hp.Helper()
	return string(res.Bytes())
}

// JSON body
func (res *ResHelper) JSON() (v interface{}) {
	res.hp.Helper()
	res.hp.E(json.Unmarshal(res.Bytes(), &v))
	return
}
