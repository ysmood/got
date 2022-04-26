package example_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/example"
)

type mockResponseWriter struct {
	got.Mock
}

func (m *mockResponseWriter) Write(b []byte) (int, error) {
	return m.Proxy("Write").(func([]byte) (int, error))(b)
}

func (m *mockResponseWriter) Header() http.Header { return nil }

func (m *mockResponseWriter) WriteHeader(c int) {}

func TestMocking(t *testing.T) {
	g := setup(t)

	m := &mockResponseWriter{}

	m.Stub("Write", func(b []byte) (int, error) {
		g.Eq(string(b), "3")
		return 0, nil
	})

	u, _ := url.Parse("?a=1&b=2")
	example.ServeSum(m, &http.Request{URL: u})

	m.On(m, "Write").When([]byte("3")).Return(1, nil)
	example.ServeSum(m, &http.Request{URL: u})
}
