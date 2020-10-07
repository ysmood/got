package got_test

import (
	"fmt"
	"testing"

	"github.com/ysmood/got"
)

var _ got.Testable = &mock{}

type mock struct {
	t      *testing.T
	failed bool
	msg    string
}

func (m *mock) Helper()        {}
func (m *mock) Cleanup(func()) {}

func (m *mock) Fail() {
	m.failed = true
}
func (m *mock) FailNow() {
	m.failed = true
	panic("fail now")
}

func (m *mock) Logf(format string, args ...interface{}) {
	m.msg = fmt.Sprintf(format, args...)
}

func (m *mock) Run(name string, fn func(*mock)) {
	fn(m)
}

func (m *mock) check(expected string) {
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
