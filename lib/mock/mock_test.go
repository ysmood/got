package mock_test

import (
	"bytes"
	"testing"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/mock"
)

type mockBuffer struct {
	mock.Mock
}

func (t *mockBuffer) Write(p []byte) (n int, err error) {
	return mock.Proxy(t, t.Write)(p)
}

func (t *mockBuffer) Len() int {
	return mock.Proxy(t, t.Len)()
}

func (t *mockBuffer) NonExists() int {
	return mock.Proxy(t, t.NonExists)()
}

func TestMock(t *testing.T) {
	g := got.T(t)

	b := bytes.NewBuffer(nil)

	m := &mockBuffer{}
	m.Fallback(b)
	mock.Stub(m, m.Write, func(p []byte) (int, error) {
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
	g.Eq(val, "you should specify the mock.Mock.Fallback")

	val = g.Panic(func() {
		m := mockBuffer{}
		m.Fallback(b)
		m.NonExists()
	})
	g.Eq(val, `*bytes.Buffer doesn't have method: NonExists`)

	g.Eq(g.Panic(func() {
		m.Stop("")
	}), "the input should be a function")
}

func TestMockUtils(t *testing.T) {
	g := got.T(t)

	b := bytes.NewBuffer(nil)

	m := &mockBuffer{}
	m.Fallback(b)

	{
		when := mock.On(m, m.Write).When([]byte{})
		when.Return(2, nil).Times(2)

		n, err := m.Write([]byte{})
		g.Nil(err)
		g.Eq(n, 2)

		n, err = m.Write([]byte{})
		g.Nil(err)
		g.Eq(n, 2)

		n, err = m.Write([]byte{})
		g.Nil(err)
		g.Eq(n, 0)

		g.Eq(when.Count(), 2)
		g.Len(m.Calls(m.Write), 3)
		g.Snapshot("calls", m.Calls(m.Write))
	}

	{
		mock.On(m, m.Write).When(mock.Any).Return(2, nil)
		n, err := m.Write([]byte{})
		g.Nil(err)
		g.Eq(n, 2)
	}

	{
		mock.On(m, m.Write).When([]byte{}).Return(2, nil).Once()

		n, err := m.Write([]byte{})
		g.Nil(err)
		g.Eq(n, 2)

		n, err = m.Write([]byte{})
		g.Nil(err)
		g.Eq(n, 0)
	}

	{
		mock.On(m, m.Write).When(true).Return(2, nil)
		v := g.Panic(func() {
			_, _ = m.Write(nil)
		})
		g.Eq(v, "No mock.StubOn.When matches: []interface {}{[]uint8(nil)}")
	}
}
