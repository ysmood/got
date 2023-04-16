package got_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ysmood/gop"
	"github.com/ysmood/got"
)

var setup = got.Setup(func(g got.G) {
	g.Parallel()
})

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

	if m.msg != "" {
		m.msg += "\n"
	}

	m.msg += fmt.Sprintf(format, args...)
}

func (m *mock) Run(_ string, fn func(*mock)) {
	fn(m)
}

func (m *mock) cleanup() {
	for _, f := range m.cleanupList {
		f()
	}
	m.cleanupList = nil
}

func (m *mock) check(expected string) {
	m.t.Helper()
	m.checkWithStyle(false, expected)
}

func (m *mock) checkWithStyle(visualizeStyle bool, expected string) {
	m.Lock()
	defer m.Unlock()

	m.t.Helper()

	if !m.failed {
		m.t.Error("should fail")
	}

	msg := ""
	if visualizeStyle {
		msg = gop.VisualizeANSI(m.msg)
	} else {
		msg = gop.StripANSI(m.msg)
	}

	if msg != expected {
		m.t.Errorf("\n\n[[[msg]]]\n\n%s\n\n[[[doesn't equal]]]\n\n%s\n\n", msg, expected)
	}

	m.failed = false
	m.msg = ""
}
