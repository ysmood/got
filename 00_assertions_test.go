package got_test

import (
	"testing"

	"github.com/ysmood/got"
)

// Use got as an light weight assertion lib in standard Go test function.

func TestAssertions(t *testing.T) {
	g := got.New(t)

	g.Eq(1, 1) // assert 1 equals 1
	g.Lt(1, 2) // assert 1 is less than 2
}
