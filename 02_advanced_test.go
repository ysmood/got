package got_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/ysmood/got"
)

// The advanced way to use got.Each

func init() {
	// Set default timeout for the entire "go test"
	got.DefaultFlags("timeout=10s")
}

func TestAdvanced(t *testing.T) {
	count := 1

	got.Each(t, func(t *testing.T) Advanced {
		g := got.New(t)

		// The function passed to it will be surely executed after the test
		g.Cleanup(func() {
			count++
		})

		// Concurrently run each test
		g.Parallel()

		// Timeout for each test
		g.PanicAfter(time.Second)

		return Advanced{g, count}
	})
}

// Usually, we use a shorter name like A or T to reduce distraction
type Advanced struct {
	// Use any assertion lib you like
	got.G

	// Share states between tests
	count int
}

func (t Advanced) A() {
	t.Desc("call t.FailNow() if 1 != 1.0").Must().Eq(1, 1.0)

	// This line won't be executed because the previous line will end the current goroutine
	t.Eq(1, 2)
}

// Use got.Skip to skip a test
func (t Advanced) B(got.Skip) {
	// Run "go doc got.Assertion" to list available assertion helpers
	t.Eq([]int{1, 2}, []int{1, 2})
}

// Use got.Only to run specific tests, same as "go test -run TestAdvanced/^C$"
func (t Advanced) C(got.Only) {
	// Run "go doc got.Utils" to list available helpers
	s := t.Serve()
	s.Route("/get", ".json", 10)

	val := t.Req("", s.URL("/get")).JSON()
	t.Eq(val, 10)

	data := map[string]interface{}{"a": 1}
	s.Mux.HandleFunc("/post", func(rw http.ResponseWriter, r *http.Request) {
		t.Eq(t.JSON(r.Body), data)
	})
	t.Req("POST", s.URL("/post"), ".json", data)
}

// Table driven tests
func (t Advanced) D() {
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

// Subtest
func (t Advanced) E() {
	t.Run("subtest", func(t got.G) {
		t.Eq(1, 1)
	})
}
