package got_test

import (
	"testing"

	"github.com/ysmood/got"
)

func TestOptions(t *testing.T) {
	m := &mock{t: t}
	g := got.NewWith(m, got.Defaults())
	g.Eq(1, 2)
	m.check("\x1b[32m1\x1b[0m\x1b[31m ⦗not ≂⦘ \x1b[0m\x1b[32m2\x1b[0m\n\n\x1b[0m\x1b[1;35m@@ -1 +1 @@\x1b[m\x1b[0m\n\x1b[0m\x1b[31m-\x1b[0m\x1b[0m\x1b[1m\x1b[37m\x1b[41m1\x1b[0m\x1b[0m\x1b[31m\x1b[0m\n\x1b[0m\x1b[32m+\x1b[0m\x1b[0m\x1b[1m\x1b[37m\x1b[42m2\x1b[0m\x1b[0m\x1b[32m\x1b[0m\n\n")

	g = got.NewWith(m, got.NoColor())
	g.Eq(1, 2)
	m.check("1 ⦗not ≂⦘ 2\n\n@@ -1 +1 @@\n-1\n+2\n\n")
}
