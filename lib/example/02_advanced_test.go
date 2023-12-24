package example_test

import (
	"testing"
	"time"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/example"
)

func TestChainMethods(t *testing.T) {
	g := setup(t)

	g.Desc("1 must equal 1").Must().Eq(example.Sum("1", "2"), "3")
}

func TestUtils(t *testing.T) {
	g := setup(t)

	// Run "go doc got.Utils" to list available helpers
	s := g.Serve()
	s.Mux.HandleFunc("/", example.ServeSum)

	val := g.Req("", s.URL("?a=1&b=2")).Bytes().String()
	g.Eq(val, "3")
}

func TestTableDriven(t *testing.T) {
	testCases := []struct{ desc, a, b, expected string }{{
		"first",
		"1", "2", "3",
	}, {
		"second",
		"2", "3", "5",
	}}

	for _, c := range testCases {
		t.Run(c.desc, func(t *testing.T) {
			g := setup(t)
			g.Eq(example.Sum(c.a, c.b), c.expected)
		})
	}
}

func TestSnapshot(t *testing.T) {
	g := setup(t)

	g.Snapshot("snapshot the map value", map[int]string{1: "1", 2: "2"})
}

func TestWaitGroup(t *testing.T) {
	g := got.T(t)

	check := func() {
		time.Sleep(time.Millisecond * 30)

		g.Eq(1, 1)
	}

	// This check won't be executed because the test will end before the goroutine starts.
	go check()

	// This check will be executed because the test will wait for the goroutine to finish.
	g.Go(check)
}
