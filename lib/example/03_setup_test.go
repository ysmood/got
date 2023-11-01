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

// G is your custom test context.
type G struct {
	got.G

	// You can add your own fields, usually data you want init before each test.
	time string
}

// setup is a helper function to setup your test context G.
var setup = func(t *testing.T) G {
	g := got.T(t)

	// The function passed to it will be surely executed after the test
	g.Cleanup(func() {})

	// Concurrently run each test
	g.Parallel()

	// Make sure there's no goroutine leak for each test
	gotrace.CheckLeak(g, 0)

	// Timeout for each test
	g.PanicAfter(time.Second)

	return G{g, time.Now().Format(time.DateTime)}
}

func TestSetup(t *testing.T) {
	g := setup(t)

	// Here we use the custom field we have defined in G
	g.Gt(g.time, "2023-01-02 15:04:05")
}
