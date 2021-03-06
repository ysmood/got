package got_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/ysmood/got"
)

// The advanced way to use got.Each

func init() {
	got.DefaultFlags("timeout=10s") // set default timeout for "go test"
}

func TestAdvanced(t *testing.T) {
	got.Each(t, setup)
}

func setup(t *testing.T) Advanced {
	opts := got.Defaults()
	opts.Keyword = func(s string) string {
		return " \x1b[31m" + s + "\x1b[0m " // print all keywords in red
	}
	g := got.NewWith(t, opts)

	g.Cleanup(func() {
		// cleanup for each test
	})

	g.Parallel() // concurrently run each test

	g.PanicAfter(time.Second) // timeout for each test

	return Advanced{g}
}

type Advanced struct { // usually, we use a shorter name like A or T to reduce distraction
	got.G // use any assertion lib you like
}

func (t Advanced) A() {
	t.Desc("call t.FailNow() if 1 != 1.0").Must().Eq(1, 1.0)
	t.Eq(1, 2) // this line won't be executed
}

func (t Advanced) B(got.Skip) { // use got.Skip to skip a test
	t.Eq([]int{1, 2}, []int{1, 2}) // run "go doc got.Assertion" to list available assertion helpers
}

func (t Advanced) C(got.Only) { // use got.Only to run specific tests, same as "go test -run TestAdvanced/^C$"
	s := t.Serve() // run "go doc got.Utils" to list available helpers
	s.Route("/get", ".json", 10)

	val := t.Req("", s.URL("/get")).JSON()
	t.Eq(val, 10)

	data := map[string]interface{}{"a": 1}
	s.Mux.HandleFunc("/post", func(rw http.ResponseWriter, r *http.Request) {
		t.Eq(t.JSON(r.Body), data)
	})
	t.Req("POST", s.URL("/post"), ".json", data)
}

func (t Advanced) D() { // table driven tests
	testCases := []struct {
		desc           string
		a, b, expected int
	}{{
		"1 + 2 = 3",
		1, 2, 3,
	}, {
		"2 + 3 = 5",
		2, 3, 5,
	}}

	add := func(a, b int) int { return a + b }

	for _, c := range testCases {
		t.Desc(c.desc).Eq(add(c.a, c.b), c.expected)
	}
}

func (t Advanced) E() { // subtests
	t.Run("subtest", func(t got.G) {
		t.Eq(1, 1)
	})
}
