package example_test

import (
	"time"

	"github.com/ysmood/got"
	"github.com/ysmood/gotrace"
)

func init() {
	// Set default timeout for the entire "go test"
	got.DefaultFlags("timeout=10s")
}

var setup = got.Setup(func(g got.G) {
	// The function passed to it will be surely executed after the test
	g.Cleanup(func() {})

	// Concurrently run each test
	g.Parallel()

	// Make sure there's no goroutine leak for each test
	gotrace.CheckLeak(g, 0)

	// Timeout for each test
	g.PanicAfter(time.Second)
})
