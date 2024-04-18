package got

import (
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Serve http on a random port. The server will be auto-closed after the test.
func (ut Utils) Serve() *Router {
	ut.Helper()
	return ut.ServeWith("tcp4", "127.0.0.1:0")
}

// ServeWith specified network and address
func (ut Utils) ServeWith(network, address string) *Router {
	ut.Helper()

	mux := http.NewServeMux()
	srv := &http.Server{Handler: mux}

	l, err := net.Listen(network, address)
	ut.err(err)

	ut.Cleanup(func() { _ = srv.Close() })

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
	p := strings.Join(path, "")
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	return rt.HostURL.String() + p
}

// Route on the pattern. Check the doc of [http.ServeMux] for the syntax of pattern.
// It will use [Utils.HandleHTTP] to handle each request.
func (rt *Router) Route(pattern, file string, value ...interface{}) *Router {
	h := rt.ut.HandleHTTP(file, value...)

	rt.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})

	return rt
}

// HandleHTTP handles a request. If file exists serve the file content. The file will be used to set the Content-Type header.
// If the file doesn't exist, the value will be encoded by [Utils.Write] and used as the response body.
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
