package got_test

import (
	"math/rand"
	"testing"

	"github.com/ysmood/got"
)

// The advanced way to use got.Each

func TestAdvanced(t *testing.T) {
	got.Each(t, setup)
}

func setup(t *testing.T) Advanced {
	t.Cleanup(func() {
		// cleanup for each test
	})

	t.Parallel() // concurrently run each test

	v := rand.Int() // generate a random value for each test

	return Advanced{got.New(t), v}
}

type Advanced struct { // usually, we use a shorter name like A or T to reduce distraction
	got.Assertion // use whatever assertion lib you like

	rand int // add any custom fields you like
}

func (t Advanced) A() {
	t.Eq(1, 1.0).Msg("b.FailNow() if %v != %v", 1, 1.0).Must()
}

func (t Advanced) B(got.Only) { // use got.Only to run specific tests
	t.Eq([]int{1, 2}, []int{1, 2})
}

func (t Advanced) C(got.Skip) { // use got.Skip to skip a test
	t.check()
}

func (t Advanced) check() { // non-exported methods won't be treated as tests
	t.Helper()

	t.Neq(t.rand, 1)
}
