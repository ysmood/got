package example_test

import (
	"testing"
	"time"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/example"
)

func TestSuite(t *testing.T) {
	// Execute each exported methods of SumSuite.
	// Each exported methods on SumSuite is a test case.
	got.Each(t, SumSuite{})
}

type SumSuite struct {
	got.G
}

func (g SumSuite) Sum() {
	g.Eq(example.Sum("1", "1"), "2")
}

func TestSumAdvancedSuite(t *testing.T) {
	// The got.Each can also accept a function to init the g for each test case.
	got.Each(t, func(t *testing.T) SumAdvancedSuite {
		g := got.New(t)

		// Concurrently run each test
		g.Parallel()

		// Timeout for each test
		g.PanicAfter(time.Second)

		return SumAdvancedSuite{g, "1", "2"}
	})
}

type SumAdvancedSuite struct {
	got.G

	a, b string
}

func (g SumAdvancedSuite) Sum() {
	g.Eq(example.Sum(g.a, g.b), "3")
}
