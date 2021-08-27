package got_test

import (
	"testing"

	"github.com/ysmood/got"
)

// Use "go test" to run the tests.

func TestSuite(t *testing.T) {
	// Execute each exported methods of Simple.
	got.Each(t, Simple{})
}

type Simple struct {
	got.G
}

// Test case A
func (t Simple) A() {
	// Assert equality of value 1 and 1.0
	t.Eq(1, 1.0)
}

// Test case B
func (t Simple) B() {
	// Assert that 1 is less than 2
	t.Lt(1, 2)
}
