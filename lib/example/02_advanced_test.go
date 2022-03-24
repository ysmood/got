package example_test

import (
	"testing"

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
