package got_test

import (
	"fmt"

	"github.com/ysmood/got"
)

// It can even run without the Go test framework
func Example_standlone() {
	tester := &T{}

	got.Each(tester, func(tt *T) t {
		return t{got.NewWith(tt, got.NoColor().NoDiff())}
	})

	// Output:
	// 1 ⦗not ≂⦘ 2
	// 1 ⦗not >⦘ 1
}

type t struct {
	got.G
}

func (t t) A() {
	t.Eq(1, 2)
}

func (t t) B() {
	t.Gt(1, 1)
}

// T is a an empty tester.
// You can config it to fit your specific requirements.
type T struct {
}

func (t *T) Run(name string, fn func(*T)) { fn(t) }

func (t *T) Logf(f string, args ...interface{}) { fmt.Printf(f+"\n", args...) }

func (t *T) Name() string { return "" }

func (t *T) Skipped() bool { return false }

func (t *T) Failed() bool { return false }

func (t *T) Helper() {}

func (t *T) Cleanup(func()) {}

func (t *T) SkipNow() {}

func (t *T) Fail() {}

func (t *T) FailNow() {}
