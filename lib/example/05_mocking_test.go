package example_test

import (
	"math/rand"
	"net/http"
	"net/url"
	"testing"

	"github.com/ysmood/got/lib/example"
	"github.com/ysmood/got/lib/mock"
)

// Mocking the http.ResponseWriter interface
type mockResponseWriter struct {
	mock.Mock
}

func (m *mockResponseWriter) Write(b []byte) (int, error) {
	// Proxy the input and output of mockResponseWriter.Write method
	return mock.Proxy(m, m.Write)(b)
}

func (m *mockResponseWriter) Header() http.Header {
	return mock.Proxy(m, m.Header)()
}

func (m *mockResponseWriter) WriteHeader(code int) {
	mock.Proxy(m, m.WriteHeader)(code)
}

func TestMocking(t *testing.T) {
	g := setup(t)

	m := &mockResponseWriter{}

	// Stub the mockResponseWriter.Write method with ours
	mock.Stub(m, m.Write, func(b []byte) (int, error) {
		// Here want to ensure the input is "3"
		g.Eq(string(b), "3")
		return 0, nil
	})

	u, _ := url.Parse("?a=1&b=2")
	example.ServeSum(m, &http.Request{URL: u})

	// When the input is "3" let the  mockResponseWriter.Write return (1, nil)
	mock.On(m, m.Write).When([]byte("3")).Return(1, nil)

	example.ServeSum(m, &http.Request{URL: u})
}

// mock the rand.Source
type mockSource struct {
	mock.Mock
}

func (m *mockSource) Int63() int64 {
	return mock.Proxy(m, m.Int63)()
}

func (m *mockSource) Seed(seed int64) {
	mock.Proxy(m, m.Seed)(seed)
}

func TestFallback(t *testing.T) {
	g := setup(t)

	m := &mockSource{}

	// Sometimes if there are a lot of method, and you want to stub only one of them,
	// you can use the Fallback method to fallback all non-stubbed methods to the struct you have passed to it.
	m.Fallback(rand.NewSource(1))

	// Here the Rand method will always return the same value.
	g.Eq(example.Rand(m), "5577006791947779410")
}
