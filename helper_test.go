package got_test

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"github.com/ysmood/got"
)

func TestHelper(t *testing.T) {
	hp := got.New(t)

	hp.Len(hp.Srand(10), 10)

	f := hp.Open(true, "tmp/test.txt")
	hp.Nil(os.Stat("tmp/test.txt"))
	hp.Write(1)(f)
	hp.Nil(f.Close())
	f = hp.Open(false, "tmp/test.txt")
	hp.Eq(hp.ReadJSON(f), 1)

	buf := bytes.NewBuffer(nil)
	hp.Write([]byte("ok"))(buf)
	hp.Eq(buf.String(), "ok")

	{
		s := hp.Serve()
		s.Route("/", ".txt")
		s.Route("/file", "go.mod")
		s.Route("/a", ".html", "ok")
		s.Route("/b", ".json", "ok", 1)
		f, err := os.Open("go.mod")
		hp.E(err)
		s.Route("/c", ".html", f)
		s.Mux.HandleFunc("/d", func(rw http.ResponseWriter, r *http.Request) {
			hp.Eq(hp.ReadString(r.Body), "1\n")
		})
		s.Mux.HandleFunc("/e", func(rw http.ResponseWriter, r *http.Request) {
			hp.Eq(hp.ReadString(r.Body), "[1,2]\n")
		})
		s.Mux.HandleFunc("/f", func(rw http.ResponseWriter, r *http.Request) {
			hp.Has(r.Header.Get("Content-Type"), "application/json")
		})

		hp.Eq(hp.Req("", s.URL()).String(), "")
		hp.Has(hp.Req("", s.URL("/file")).String(), "ysmood/got")
		hp.Eq(hp.Req("", s.URL("/a")).String(), "ok")
		hp.Eq(hp.Req("", s.URL("/a")).String(), "ok")
		res := hp.Req("", s.URL("/b"))
		hp.Eq(res.JSON(), []interface{}{"ok", 1})
		hp.Has(res.Header.Get("Content-Type"), "application/json")
		hp.Has(hp.Req("", s.URL("/c")).String(), "ysmood/got")
		hp.Req(http.MethodPost, s.URL("/d"), 1)
		hp.Req(http.MethodPost, s.URL("/e"), "", 1, 2)
		hp.Req(http.MethodPost, s.URL("/f"), ".json", 1)
	}
}
