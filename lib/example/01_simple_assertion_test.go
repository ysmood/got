package example_test

import (
	"testing"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/example"
)

// Use got as an light weight assertion lib in standard Go test function.
func TestAssertion(t *testing.T) {
	// Run "go doc got.Assertions" to list available assertion methods.
	got.T(t).Eq(example.Sum("1", "1"), "2")
}
