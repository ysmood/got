package example_test

import (
	"fmt"
	"testing"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/diff"
	"github.com/ysmood/got/lib/gop"
)

// An example to only override the default error output of got.Assertions.Eq
func TestCustomizeAssertionOutput(t *testing.T) {
	g := got.New(t)

	dh := got.NewDefaultAssertionError(gop.ThemeDefault, diff.ThemeDefault)
	h := got.AssertionErrorReport(func(c *got.AssertionCtx) string {
		if c.Type == got.AssertionEq {
			return fmt.Sprintf("%v != %v", c.Details[0], c.Details[1])
		}
		return dh.Report(c)
	})
	g.Assertions.ErrorHandler = h

	g.Eq(1, 1)
}
