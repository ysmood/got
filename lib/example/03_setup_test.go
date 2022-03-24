package example_test

import (
	"testing"
	"time"

	"github.com/ysmood/got"
	"github.com/ysmood/gotrace"
)

func init() {
	// Set default timeout for the entire "go test"
	got.DefaultFlags("timeout=10s")
}

func TestMain(m *testing.M) {
	// Make sure we don't leaking goroutines
	gotrace.CheckMain(m, 0)
}

var setup = got.Setup(func(g got.G) {
	// The function passed to it will be surely executed after the test
	g.Cleanup(func() {})

	if got.Parallel() > 0 {
		// Concurrently run each test
		g.Parallel()
	} else {
		// Make sure there's no goroutine leak for each test
		gotrace.Check(g, 0)
	}

	// Timeout for each test
	g.PanicAfter(time.Second)
})
