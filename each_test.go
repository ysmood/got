package got_test

import (
	"testing"

	"github.com/ysmood/got"
)

func TestEach(t *testing.T) {
	got.Each(t, StructVal{val: 1})
}

type StructVal struct {
	got.Assertion
	val int
}

func (c StructVal) A() {
	c.Eq(c.val, 1)
}

func TestEachSkip(t *testing.T) {
	got.Each(t, Container{})
}

type Container struct {
	Embedded
}

func (c Container) A() { c.Fail() }
func (c Container) B() {}

type Embedded struct {
	*testing.T
}

func (c Embedded) A() {}
func (c Embedded) C() { c.Fail() }

func TestEachErr(t *testing.T) {
	as := got.New(t)
	m := &mock{t: t}

	as.Panic(func() {
		got.Each(m, nil)
	})
	m.check("[got.Each] iteratee shouldn't be nil")

	as.Panic(func() {
		got.Each(m, 1)
	})
	m.check("[got.Each] iteratee <int> should be a struct or <func(got.Testable) Ctx>")

	it := func() Err { return Err{} }
	as.Panic(func() {
		got.Each(m, it)
	})
	m.check("[got.Each] iteratee <func() got_test.Err> should be a struct or <func(got.Testable) Ctx>")

	as.Panic(func() {
		got.Each(m, func(t *mock) Err {
			return Err{}
		})
	})
	m.check("[got.Each] got_test.Err.A shouldn't have arguments or return values")
}

type Err struct {
}

func (s Err) A(int) {}

func TestParallelErr(t *testing.T) {
	ok := testing.RunTests(
		func(pat, str string) (bool, error) { return true, nil },
		[]testing.InternalTest{{
			F: func(t *testing.T) {
				got.Each(t, func(t *testing.T) ParallelErr {
					t.Parallel()
					return ParallelErr{t}
				})
			},
		}},
	)
	got.New(t).False(ok)
}

type ParallelErr struct {
	*testing.T
}

func (p ParallelErr) A() {
	p.Fail()
}

func (p ParallelErr) B() {
}
