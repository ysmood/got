package got_test

import (
	"bytes"
	"testing"

	"github.com/ysmood/got"
)

type mockBuffer struct {
	got.Mock
}

func (t *mockBuffer) Write(p []byte) (n int, err error) {
	return t.Proxy("Write").(func([]byte) (int, error))(p)
}

func (t *mockBuffer) Len() int {
	return t.Proxy("Len").(func() int)()
}

func (t *mockBuffer) Nonexists() int {
	return t.Proxy("Nonexists").(func() int)()
}

func TestMock(t *testing.T) {
	g := setup(t)

	b := bytes.NewBuffer(nil)

	m := mockBuffer{}
	m.Fallback(b)
	m.Stub("Write", func(p []byte) (int, error) {
		return b.Write(append(p, []byte("  ")...))
	})
	n, err := m.Write([]byte("test"))
	g.Nil(err)
	g.Eq(n, 6)

	g.Eq(m.Len(), 6)

	val := g.Panic(func() {
		m := mockBuffer{}
		m.Len()
	})
	g.Eq(val, "you should specify the got.Mock.Origin")

	val = g.Panic(func() {
		m := mockBuffer{}
		m.Fallback(b)
		m.Nonexists()
	})
	g.Eq(val, `*bytes.Buffer doesn't have method: Nonexists`)
}

func TestMockUtils(t *testing.T) {
	g := setup(t)

	b := bytes.NewBuffer(nil)

	m := &mockBuffer{}
	m.Fallback(b)

	{
		m.On(m, "Write").When([]byte{}).Return(2, nil).Times(2)

		n, err := m.Write(nil)
		g.Nil(err)
		g.Eq(n, 2)

		n, err = m.Write(nil)
		g.Nil(err)
		g.Eq(n, 2)

		n, err = m.Write(nil)
		g.Nil(err)
		g.Eq(n, 0)
	}

	{
		m.On(m, "Write").When(got.Any).Return(2, nil)
		n, err := m.Write(nil)
		g.Nil(err)
		g.Eq(n, 2)
	}

	{
		m.On(m, "Write").When([]byte{}).Return(2, nil).Once()

		n, err := m.Write(nil)
		g.Nil(err)
		g.Eq(n, 2)

		n, err = m.Write(nil)
		g.Nil(err)
		g.Eq(n, 0)
	}

	{
		m.On(m, "Write").When(true).Return(2, nil)
		v := g.Panic(func() {
			_, _ = m.Write(nil)
		})
		g.Eq(v, "No got.StubOn.When matches: []interface {}{[]uint8(nil)}")
	}
}
