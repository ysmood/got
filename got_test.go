package got_test

import (
	"testing"

	"github.com/ysmood/got"
)

func TestSetup(t *testing.T) {
	g := setup(t)
	g.Eq(1, 1)
}

func TestOptions(t *testing.T) {
	m := &mock{t: t}
	g := got.NewWith(m, got.Defaults())
	g.Eq(1, 2)
	m.check("\x1b[32m1\x1b[0m\x1b[31m ⦗not ==⦘ \x1b[0m\x1b[32m2\x1b[0m\n\n1   \x1b[31m- \x1b[0m\x1b[31m1\n\x1b[0m  1 \x1b[32m+ \x1b[0m\x1b[32m2\n\x1b[0m\n")

	g = got.NewWith(m, got.NoColor())
	g.Eq(1, 2)
	m.check("1 ⦗not ==⦘ 2\n\n1   \x1b[31m- \x1b[0m\x1b[31m1\n\x1b[0m  1 \x1b[32m+ \x1b[0m\x1b[32m2\n\x1b[0m\n")
}
