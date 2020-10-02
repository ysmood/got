package got_test

import (
	"math/rand"
	"testing"

	"github.com/ysmood/got"
)

func TestBasic(t *testing.T) {
	got.Each(t, setup)
}

func setup(t *testing.T) Basic {
	t.Cleanup(func() {
		// cleanup ...
	})

	t.Parallel() // concurrently run each test

	v := rand.Int() // generate a random value for each test

	return Basic{got.New(t), v}
}

type Basic struct {
	got.Assertion // use whatever assertion lib you like

	rand int // add any custom fields you like
}

func (b Basic) A() {
	b.Eq(1, 1.0).Msg("b.FailNow() if %v != %v", 1, 1.0).Must()
}

func (b Basic) B() {
	b.Eq([]int{1, 2}, []int{1, 2})
}

func (b Basic) C() {
	b.check()
}

func (b Basic) check() { // non-exported methods won't be treated as tests
	b.Helper()

	b.Neq(b.rand, 1)
}
