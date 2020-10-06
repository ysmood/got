package got_test

import (
	"testing"

	"github.com/ysmood/got"
)

func TestEach(t *testing.T) {
	count := got.Each(t, StructVal{val: 1})
	got.New(t).Eq(count, 2)
}

type StructVal struct {
	got.Assertion
	val int
}

func (c StructVal) Normal() {
	c.Eq(c.val, 1)
}

func (c StructVal) ExtraInOut(int) int {
	c.Eq(c.val, 1)
	return 0
}

func TestEachEmbedded(t *testing.T) {
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
	m.check("iteratee shouldn't be nil")

	as.Panic(func() {
		got.Each(m, 1)
	})
	m.check("iteratee <int> should be a struct or <func(got.Testable) Ctx>")

	it := func() Err { return Err{} }
	as.Panic(func() {
		got.Each(m, it)
	})
	m.check("iteratee <func() got_test.Err> should be a struct or <func(got.Testable) Ctx>")
}

type Err struct {
}

func (s Err) A(int) {}
