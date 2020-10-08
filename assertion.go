package got

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

// Result helper
type Result struct {
	as     G
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

// Eq a ≂ b
func (as G) Eq(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) == 0 {
		return
	}
	return as.err("%s %s %s", as.d(a), as.k("not ≂"), as.d(b))
}

// Neq a != b
func (as G) Neq(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) != 0 {
		return
	}
	return as.err("%s %s %s", as.d(a), as.k("not ≠"), as.d(b))
}

// Equal a == b
func (as G) Equal(a, b interface{}) (result Result) {
	as.Helper()
	if a == b {
		return
	}
	return as.err("%s %s %s", as.d(a), as.k("not =="), as.d(b))
}

// Gt a > b
func (as G) Gt(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) > 0 {
		return
	}
	return as.err("%s %s %s", as.d(a), as.k("not >"), as.d(b))
}

// Gte a >= b
func (as G) Gte(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) >= 0 {
		return
	}
	return as.err("%s %s %s", as.d(a), as.k("not ≥"), as.d(b))
}

// Lt a < b
func (as G) Lt(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) < 0 {
		return
	}
	return as.err("%s %s %s", as.d(a), as.k("not <"), as.d(b))
}

// Lte a <= b
func (as G) Lte(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) <= 0 {
		return
	}
	return as.err("%s %s %s", as.d(a), as.k("not ≤"), as.d(b))
}

// True a == true
func (as G) True(a bool) (result Result) {
	as.Helper()
	if a {
		return
	}
	return as.err("%s", as.k("should be <true>"))
}

// False a == false
func (as G) False(a bool) (result Result) {
	as.Helper()
	if !a {
		return
	}
	return as.err("%s", as.k("should be <false>"))
}

// Nil fails if last arg is not nil
func (as G) Nil(args ...interface{}) (result Result) {
	as.Helper()
	if len(args) == 0 {
		return as.err("%s", as.k("no args received"))
	}
	last := args[len(args)-1]
	if isNil(last) {
		return
	}
	return as.err("%s %s %s", as.k("last value"), as.d(last), as.k("should be <nil>"))
}

// NotNil fails if last arg is nil
func (as G) NotNil(args ...interface{}) (result Result) {
	as.Helper()
	if len(args) == 0 {
		return as.err("%s", as.k("no args received"))
	}
	last := args[len(args)-1]
	if !isNil(last) {
		return
	}
	if last == nil {
		return as.err("%s", as.k("last value shouldn't be <nil>"))
	}
	return as.err("<%s> %s", reflect.TypeOf(last), as.k("shouldn't be <nil>"))
}

// Regex matches str
func (as G) Regex(pattern, str string) (result Result) {
	as.Helper()
	if regexp.MustCompile(pattern).MatchString(str) {
		return
	}
	return as.err("%s %s %s", pattern, as.k("should match"), str)
}

// Has str in container
func (as G) Has(container, str string) (result Result) {
	as.Helper()
	if strings.Contains(container, str) {
		return
	}
	return as.err("%s %s %s", container, as.k("should has"), str)
}

// Len len(list) == l
func (as G) Len(list interface{}, l int) (result Result) {
	as.Helper()
	actual := reflect.ValueOf(list).Len()
	if actual == l {
		return
	}
	return as.err("%s %d %s %d", as.k("expect len"), actual, as.k("to be"), l)
}

// Err fails if last arg is not error
func (as G) Err(args ...interface{}) (result Result) {
	as.Helper()
	if len(args) == 0 {
		return as.err("%s", as.k("no args received"))
	}
	last := args[len(args)-1]
	if err, _ := last.(error); err != nil {
		return
	}
	return as.err("%s %s %s", as.k("last value"), as.d(last), as.k("should be <error>"))
}

// E is a shortcut for Nil(args...).Must()
func (as G) E(args ...interface{}) {
	as.Helper()
	as.Nil(args...).Must()
}

// Panic fails if fn doesn't panic
func (as G) Panic(fn func()) (result Result) {
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

// Is fails if a is not kind of b
func (as G) Is(a, b interface{}) (result Result) {
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
