package got_test

import (
	"testing"
)

func TestSetup(t *testing.T) {
	g := setup(t)
	g.Eq(1, 1)
}
