package got

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Assertion is the assertion context
type Assertion struct {
	Testable

	d func(v interface{}) string
	k func(string) string
}

// Options for Assertion
type Options struct {
	// Dump a value to human readable string
	Dump func(interface{}) string

	// Format keywords in the assertion message.
	// Such as color it for CLI output.
	Keyword func(string) string
}

// Defaults for Options
func Defaults() Options {
	return Options{
		func(v interface{}) string {
			return fmt.Sprintf("%v (%v)", v, reflect.TypeOf(v))
		},
		func(s string) string {
			return "⦗" + s + "⦘"
		},
	}
}

// New assertion helper
func New(t Testable) Assertion {
	return NewWith(t, Defaults())
}

// NewWith assertion helper with options
func NewWith(t Testable, opts Options) Assertion {
	return Assertion{t, opts.Dump, opts.Keyword}
}

// Result helper
type Result struct {
	as     Assertion
	failed bool
}

// Msg if fails
func (r Result) Msg(format string, args ...interface{}) Result {
	if r.failed {
		r.as.Helper()
		r.as.Logf(format, args...)
	}
	return r
}

// Must FailNow if fails
func (r Result) Must() Result {
	if r.failed {
		r.as.Helper()
		r.as.FailNow()
	}
	return r
}

// Eq a == b
func (as Assertion) Eq(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) == 0 {
		return
	}
	return as.err("%s %s %s", as.d(a), as.k("≂"), as.d(b))
}

// Neq a != b
func (as Assertion) Neq(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) != 0 {
		return
	}
	return as.err("%s %s %s", as.d(a), as.k("≠"), as.d(b))
}

// Gt a > b
func (as Assertion) Gt(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) > 0 {
		return
	}
	return as.err("%s %s %s", as.d(a), as.k(">"), as.d(b))
}

// Gte a >= b
func (as Assertion) Gte(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) >= 0 {
		return
	}
	return as.err("%s %s %s", as.d(a), as.k("≥"), as.d(b))
}

// Lt a < b
func (as Assertion) Lt(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) < 0 {
		return
	}
	return as.err("%s %s %s", as.d(a), as.k("<"), as.d(b))
}

// Lte a <= b
func (as Assertion) Lte(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) <= 0 {
		return
	}
	return as.err("%s %s %s", as.d(a), as.k("≤"), as.d(b))
}

// True a == true
func (as Assertion) True(a bool) (result Result) {
	as.Helper()
	if a {
		return
	}
	return as.err("%s", as.k("should be <true>"))
}

// False a == false
func (as Assertion) False(a bool) (result Result) {
	as.Helper()
	if !a {
		return
	}
	return as.err("%s", as.k("should be <false>"))
}

// Nil args[-1] == nil
func (as Assertion) Nil(args ...interface{}) (result Result) {
	as.Helper()
	if len(args) == 0 {
		return as.err("%s", as.k("no args received"))
	}
	last := args[len(args)-1]
	if isNil(last) {
		return
	}
	return as.err("%s %s", as.d(last), as.k("should be <nil>"))
}

// NotNil args[-1] != nil
func (as Assertion) NotNil(args ...interface{}) (result Result) {
	as.Helper()
	if len(args) == 0 {
		return as.err("%s", as.k("no args received"))
	}
	last := args[len(args)-1]
	if !isNil(last) {
		return
	}
	return as.err("%s %s", as.d(last), as.k("shouldn't be <nil>"))
}

// Regex matches str
func (as Assertion) Regex(pattern, str string) (result Result) {
	as.Helper()
	if regexp.MustCompile(pattern).MatchString(str) {
		return
	}
	return as.err("%s %s %s", pattern, as.k("should match"), str)
}

// Has str in container
func (as Assertion) Has(container, str string) (result Result) {
	as.Helper()
	if strings.Contains(container, str) {
		return
	}
	return as.err("%s %s %s", container, as.k("should has"), str)
}

// Len len(list) == l
func (as Assertion) Len(list interface{}, l int) (result Result) {
	as.Helper()
	actual := reflect.ValueOf(list).Len()
	if actual == l {
		return
	}
	return as.err("%s %d %s %d", as.k("expect len"), actual, as.k("to be"), l)
}

// Err args[-1] is error and not nil
func (as Assertion) Err(args ...interface{}) (result Result) {
	as.Helper()
	if len(args) == 0 {
		return as.err("%s", as.k("no args received"))
	}
	last := args[len(args)-1]
	if err, _ := last.(error); err != nil {
		return
	}
	return as.err("%s %s", as.d(last), as.k("should be <error>"))
}

// Panic fn should panic
func (as Assertion) Panic(fn func()) (result Result) {
	as.Helper()

	defer func() {
		as.Helper()

		val := recover()
		if val == nil {
			result = as.err("%s", as.k("should panic"))
		}
	}()

	fn()

	return
}

// Is a a kind of b
func (as Assertion) Is(a, b interface{}) (result Result) {
	as.Helper()

	if ae, ok := a.(error); ok {
		if be, ok := b.(error); ok {
			if ae == be {
				return
			}

			if errors.Is(ae, be) {
				return
			}
			return as.err("%s %s %s", as.d(a), as.k("should in chain of"), as.d(b))
		}
	}

	at := reflect.TypeOf(a)
	bt := reflect.TypeOf(b)
	if at.Kind() == bt.Kind() {
		return
	}
	return as.err("%s %s %s", as.d(a), as.k("should kind of"), as.d(b))
}
