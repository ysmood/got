package got_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ysmood/got"
)

var _ got.Testable = &mock{}

type mock struct {
	sync.Mutex
	t           *testing.T
	failed      bool
	skipped     bool
	msg         string
	cleanupList []func()
	recover     bool
}

func (m *mock) Name() string     { return "mock" }
func (m *mock) Skipped() bool    { return m.skipped }
func (m *mock) Failed() bool     { return m.failed }
func (m *mock) Helper()          {}
func (m *mock) Cleanup(f func()) { m.cleanupList = append([]func(){f}, m.cleanupList...) }
func (m *mock) SkipNow()         {}
func (m *mock) Fail()            { m.failed = true }

func (m *mock) FailNow() {
	m.Lock()
	defer m.Unlock()

	m.failed = true
	if !m.recover {
		panic("fail now")
	}
	m.recover = false
}

func (m *mock) Logf(format string, args ...interface{}) {
	m.Lock()
	defer m.Unlock()

	m.msg = fmt.Sprintf(format, args...)
}

func (m *mock) Run(name string, fn func(*mock)) {
	fn(m)
}

func (m *mock) cleanup() {
	for _, f := range m.cleanupList {
		f()
	}
	m.cleanupList = nil
}

func (m *mock) check(expected string) {
	m.Lock()
	defer m.Unlock()

	m.t.Helper()

	as := got.NewWith(m.t, got.Options{
		Dump: func(i interface{}) string {
			return fmt.Sprintf("\n%v\n", i)
		},
		Keyword: func(s string) string {
			return s
		},
	})

	as.True(m.failed)
	as.Eq(m.msg, expected)

	m.failed = false
	m.msg = ""
}
