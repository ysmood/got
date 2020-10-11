package got_test

import (
	"testing"

	"github.com/ysmood/got"
)

// The simplest form to use got.Each .
// No magic, use "go test" to run the tests.

func TestSimple(t *testing.T) { // entry point is just a normal Go test function
	got.Each(t, Simple{})
}

type Simple struct { // all exported methods on it will be executed as the subtests of TestSimple
}

func (t Simple) A() { // a test
}

func (t Simple) B() { // another test
}
