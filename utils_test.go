package got_test

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/ysmood/got"
)

func init() {
	got.DefaultFlags("parallel=3")
}

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

	ut.Eq(ut.ToJSONString(10), "10")

	buf := bytes.NewBuffer(nil)
	ut.Write([]byte("ok"))(buf)
	ut.Eq(buf.String(), "ok")

	ut.Run("subtest", func(t got.G) {
		t.Eq(1, 1)
	})

	ut.Eq(got.Parallel(), 3)

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
			ut.Eq(ut.Read(r.Body).String(), "1\n")
		})
		s.Mux.HandleFunc("/f", func(rw http.ResponseWriter, r *http.Request) {
			ut.Has(r.Header.Get("Content-Type"), "application/json")
			ut.Eq(r.Header.Get("Test-Header"), "ok")
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
		ut.Req(http.MethodPost, s.URL("/f"), http.Header{"Test-Header": {"ok"}}, got.ReqMIME(".json"), 1)
	}

	m := &mock{t: t}
	mut := got.New(m)

	m.msg = ""
	mut.Log("a", 1)
	ut.Eq(m.msg, "a 1\n")

	m.msg = ""
	ut.Panic(func() {
		buf := bytes.NewBufferString("a")
		mut.JSON(buf)
	})
	ut.Eq(m.msg, "invalid character 'a' looking for beginning of value\n")

	m.msg = ""
	ut.Panic(func() {
		mut.Fatal("test skip")
	})
	ut.Eq(m.msg, "test skip\n")

	m.msg = ""
	ut.Panic(func() {
		mut.Fatalf("test skip")
	})
	ut.Eq(m.msg, "test skip")

	m.msg = ""
	mut.Error("test skip")
	ut.Eq(m.msg, "test skip\n")

	m.msg = ""
	mut.Errorf("test skip")
	ut.Eq(m.msg, "test skip")

	m.msg = ""
	mut.Skip("test skip")
	ut.Eq(m.msg, "test skip\n")

	m.msg = ""
	mut.Skipf("test skip")
	ut.Eq(m.msg, "test skip")
}

func TestServe(t *testing.T) {
	ut := got.New(t)

	key := ut.Srand(8)
	s := ut.Serve().Route("/", "", key)
	count := 30

	wg := sync.WaitGroup{}
	wg.Add(count)

	request := func() {
		req, err := http.NewRequest(http.MethodGet, s.URL(), nil)
		ut.E(err)

		res, err := http.DefaultClient.Do(req)
		ut.E(err)

		b, err := io.ReadAll(res.Body)
		ut.E(err)

		ut.Eq(string(b), key)
		wg.Done()
	}

	for i := 0; i < count; i++ {
		go request()
	}

	wg.Wait()
}
