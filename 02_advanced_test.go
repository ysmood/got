package got_test

import (
	"net/http"
	"testing"

	"github.com/ysmood/got"
)

// The advanced way to use got.Each

func TestAdvanced(t *testing.T) {
	got.Each(t, setup)
}

func setup(t *testing.T) Advanced {
	t.Cleanup(func() {
		// cleanup for each test
	})

	t.Parallel() // concurrently run each test

	return Advanced{got.New(t)}
}

type Advanced struct { // usually, we use a shorter name like A or T to reduce distraction
	got.G // use any assertion lib you like
}

func (t Advanced) A() {
	t.Eq(1, 1.0).Msg("b.FailNow() if %v != %v", 1, 1.0).Must()
}

func (t Advanced) B(got.Skip) { // use got.Skip to skip a test
	t.Eq([]int{1, 2}, []int{1, 2})
}

func (t Advanced) C(got.Only) { // use got.Only to run specific tests
	s := t.Serve()
	s.Route("/get", ".json", 10)

	val := t.Req("", s.URL("/get")).JSON()
	t.Eq(val, 10)

	data := map[int]string{1: "ok"}
	s.Mux.HandleFunc("/post", func(rw http.ResponseWriter, r *http.Request) {
		t.Eq(t.ReadJSON(r.Body), data)
	})
	t.Req("POST", s.URL("/post"), ".json", data)
}
