package got_test

import (
	"bytes"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ysmood/got"
)

func TestHelper(t *testing.T) {
	ut := got.New(t)

	ctx := ut.Context()
	ctx.Cancel()
	<-ctx.Done()
	<-ut.Timeout(0).Done()

	ut.Len(ut.Srand(10), 10)

	f := ut.Open(true, "tmp/test.txt")
	ut.Nil(os.Stat("tmp/test.txt"))
	ut.Write(1)(f)
	ut.Nil(f.Close())
	f = ut.Open(false, "tmp/test.txt")
	ut.Eq(ut.JSON(f), 1)

	ut.Eq(ut.JSON([]byte("1")), 1)
	ut.Eq(ut.JSON("true"), true)

	buf := bytes.NewBuffer(nil)
	ut.Write([]byte("ok"))(buf)
	ut.Eq(buf.String(), "ok")

	{
		s := ut.Serve()
		s.Route("/", ".txt")
		s.Route("/file", "go.mod")
		s.Route("/a", ".html", "ok")
		s.Route("/b", ".json", "ok", 1)
		f, err := os.Open("go.mod")
		ut.E(err)
		s.Route("/c", ".html", f)
		s.Mux.HandleFunc("/d", func(rw http.ResponseWriter, r *http.Request) {
			ut.Eq(ut.ReadString(r.Body), "1\n")
		})
		s.Mux.HandleFunc("/e", func(rw http.ResponseWriter, r *http.Request) {
			ut.Eq(ut.ReadString(r.Body), "[1,2]\n")
		})
		s.Mux.HandleFunc("/f", func(rw http.ResponseWriter, r *http.Request) {
			ut.Has(r.Header.Get("Content-Type"), "application/json")
		})

		ut.Eq(ut.Req("", s.URL()).String(), "")
		ut.Has(ut.Req("", s.URL("/file")).String(), "ysmood/got")
		ut.Eq(ut.Req("", s.URL("/a")).String(), "ok")
		ut.Eq(ut.Req("", s.URL("/a")).String(), "ok")
		res := ut.Req("", s.URL("/b"))
		ut.Eq(res.JSON(), []interface{}{"ok", 1})
		ut.Has(res.Header.Get("Content-Type"), "application/json")
		ut.Has(ut.Req("", s.URL("/c")).String(), "ysmood/got")
		ut.Req(http.MethodPost, s.URL("/d"), 1)
		ut.Req(http.MethodPost, s.URL("/e"), "", 1, 2)
		ut.Req(http.MethodPost, s.URL("/f"), ".json", 1)
	}

	m := &mock{t: t}
	mut := got.New(m)

	mut.Log("a", 1)
	ut.Eq(m.msg, "a 1\n")

	ut.Panic(func() {
		buf := bytes.NewBufferString("a")
		mut.JSON(buf)
	})
	ut.Eq(m.msg, "invalid character 'a' looking for beginning of value\n")

	ut.Panic(func() {
		mut.Fatal("test skip")
	})
	ut.Eq(m.msg, "test skip\n")

	ut.Panic(func() {
		mut.Fatalf("test skip")
	})
	ut.Eq(m.msg, "test skip")

	mut.Error("test skip")
	ut.Eq(m.msg, "test skip\n")

	mut.Errorf("test skip")
	ut.Eq(m.msg, "test skip")

	mut.Skip("test skip")
	ut.Eq(m.msg, "test skip\n")

	mut.Skipf("test skip")
	ut.Eq(m.msg, "test skip")

	mut.HeartBeat(time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	m.cleanup()
	m.Lock()
	ut.Eq(m.msg, "[got.Utils.HeartBeat] continuing...")
	m.Unlock()
}
