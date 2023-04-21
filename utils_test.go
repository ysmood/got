package got_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ysmood/got"
)

func init() {
	got.DefaultFlags("parallel=3")
}

func TestHelper(t *testing.T) {
	ut := got.T(t)

	ctx := ut.Context()
	ctx.Cancel()
	<-ctx.Done()
	<-ut.Timeout(0).Done()

	ut.Len(ut.RandStr(10), 10)
	ut.Lt(ut.RandInt(0, 1), 1)
	ut.Gt(ut.RandInt(-2, -1), -3)

	f := ut.Open(true, "tmp/test.txt")
	ut.Nil(os.Stat("tmp/test.txt"))
	ut.Write(1)(f)
	ut.Nil(f.Close())
	f = ut.Open(false, "tmp/test.txt")
	ut.Eq(ut.JSON(f), 1)

	ut.Setenv(ut.RandStr(8), ut.RandStr(8))
	ut.MkdirAll(0, "tmp/a/b/c")

	s := ut.RandStr(16)
	ut.WriteFile("tmp/test.txt", s)
	ut.Eq(ut.Read("tmp/test.txt").String(), s)
	ut.Eq(ut.Read(123).String(), "123")
	ut.Eq(ut.Read([]byte("ok")).String(), "ok")
	ut.Eq(ut.Render("{{.}}", 10).String(), "10")

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
		ut.Eq(res.JSON(), []interface{}{"ok", float64(1)})
		ut.Has(res.Header.Get("Content-Type"), "application/json")
		ut.Has(ut.Req("", s.URL("/c")).String(), "ysmood/got")
		ut.Req(http.MethodPost, s.URL("/d"), 1)
		ut.Req(http.MethodPost, s.URL("/f"), http.Header{"Test-Header": {"ok"}}, got.ReqMIME(".json"), 1)
	}

	ut.DoAfter(time.Hour, func() {})

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
	ut := setup(t)

	key := ut.RandStr(8)
	s := ut.Serve().Route("/", "", key)
	count := 30

	wg := sync.WaitGroup{}
	wg.Add(count)

	request := func() {
		req, err := http.NewRequest(http.MethodGet, s.URL(), nil)
		ut.E(err)

		res, err := http.DefaultClient.Do(req)
		ut.E(err)

		b, err := ioutil.ReadAll(res.Body)
		ut.E(err)

		ut.Eq(string(b), key)
		wg.Done()
	}

	for i := 0; i < count; i++ {
		go request()
	}

	wg.Wait()
}

func TestPathExists(t *testing.T) {
	g := got.T(t)

	g.False(g.PathExists("not-exists"))
	g.False(g.PathExists("*!"))
	g.True(g.PathExists("lib"))
}

func TestChdir(t *testing.T) {
	g := got.T(t)

	g.Chdir("lib")

	g.PathExists("diff")
}
