package got

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Result helper
type Result struct {
	as     Assertions
	failed bool
}

// Assertions helpers
type Assertions struct {
	Testable

	d  func(v interface{}) string    // Options.Dump
	k  func(string) string           // Options.Keyword
	df func(a, b interface{}) string // Options.Diff
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
func (as Assertions) Eq(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) == 0 {
		return
	}
	return as.err("%s%s%s%s", as.d(a), as.k("not ≂"), as.d(b), as.diff(a, b))
}

// Neq a != b
func (as Assertions) Neq(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) != 0 {
		return
	}
	return as.err("%s%s%s", as.d(a), as.k("not ≠"), as.d(b))
}

// Equal a == b
func (as Assertions) Equal(a, b interface{}) (result Result) {
	as.Helper()
	if a == b {
		return
	}
	return as.err("%s%s%s%s", as.d(a), as.k("not =="), as.d(b), as.diff(a, b))
}

// Gt a > b
func (as Assertions) Gt(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) > 0 {
		return
	}
	return as.err("%s%s%s", as.d(a), as.k("not >"), as.d(b))
}

// Gte a >= b
func (as Assertions) Gte(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) >= 0 {
		return
	}
	return as.err("%s%s%s", as.d(a), as.k("not ≥"), as.d(b))
}

// Lt a < b
func (as Assertions) Lt(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) < 0 {
		return
	}
	return as.err("%s%s%s", as.d(a), as.k("not <"), as.d(b))
}

// Lte a <= b
func (as Assertions) Lte(a, b interface{}) (result Result) {
	as.Helper()
	if compare(a, b) <= 0 {
		return
	}
	return as.err("%s%s%s", as.d(a), as.k("not ≤"), as.d(b))
}

// True a == true
func (as Assertions) True(a bool) (result Result) {
	as.Helper()
	if a {
		return
	}
	return as.err("%s", as.k("should be <true>"))
}

// False a == false
func (as Assertions) False(a bool) (result Result) {
	as.Helper()
	if !a {
		return
	}
	return as.err("%s", as.k("should be <false>"))
}

// Nil fails if last arg is not nil
func (as Assertions) Nil(args ...interface{}) (result Result) {
	as.Helper()
	if len(args) == 0 {
		return as.err("%s", as.k("no args received"))
	}
	last := args[len(args)-1]
	if isNil(last) {
		return
	}
	return as.err("%s%s%s", as.k("last value"), as.d(last), as.k("should be <nil>"))
}

// NotNil fails if last arg is nil
func (as Assertions) NotNil(args ...interface{}) (result Result) {
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
	return as.err("<%s>%s", reflect.TypeOf(last), as.k("shouldn't be <nil>"))
}

// Regex matches str
func (as Assertions) Regex(pattern, str string) (result Result) {
	as.Helper()
	if regexp.MustCompile(pattern).MatchString(str) {
		return
	}
	return as.err("%s%s%s", pattern, as.k("should match"), str)
}

// Has str in container
func (as Assertions) Has(container, str string) (result Result) {
	as.Helper()
	if strings.Contains(container, str) {
		return
	}
	return as.err("%s%s%s", container, as.k("should has"), str)
}

// Len len(list) == l
func (as Assertions) Len(list interface{}, l int) (result Result) {
	as.Helper()
	actual := reflect.ValueOf(list).Len()
	if actual == l {
		return
	}
	return as.err("%s%d%s%d", as.k("expect len"), actual, as.k("to be"), l)
}

// Err fails if last arg is not error
func (as Assertions) Err(args ...interface{}) (result Result) {
	as.Helper()
	if len(args) == 0 {
		return as.err("%s", as.k("no args received"))
	}
	last := args[len(args)-1]
	if err, _ := last.(error); err != nil {
		return
	}
	return as.err("%s%s%s", as.k("last value"), as.d(last), as.k("should be <error>"))
}

// E is a shortcut for Nil(args...).Must()
func (as Assertions) E(args ...interface{}) {
	as.Helper()
	as.Nil(args...).Must()
}

// Panic fails if fn doesn't panic
func (as Assertions) Panic(fn func()) (result Result) {
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
func (as Assertions) Is(a, b interface{}) (result Result) {
	as.Helper()

	if a == nil && b == nil {
		return
	}

	if ae, ok := a.(error); ok {
		if be, ok := b.(error); ok {
			if ae == be {
				return
			}

			if errors.Is(ae, be) {
				return
			}
			return as.err("%s%s%s", as.d(a), as.k("should in chain of"), as.d(b))
		}
	}

	at := reflect.TypeOf(a)
	bt := reflect.TypeOf(b)
	if a != nil && b != nil && at.Kind() == bt.Kind() {
		return
	}
	return as.err("%s%s%s", as.d(a), as.k("should be kind of"), as.d(b))
}

// Count returns a function that must be called with the specified times
func (as Assertions) Count(times int) func() {
	as.Helper()
	count := 0

	as.Cleanup(func() {
		if count != times {
			as.Helper()
			as.Logf("Should count %d times, but got %d", times, count)
			as.Fail()
		}
	})

	return func() {
		count++
	}
}

func (as Assertions) err(format string, args ...interface{}) Result {
	as.Helper()
	as.Logf(format, args...)
	as.Fail()
	return Result{as, true}
}

func castType(a, b interface{}) interface{} {
	ta := reflect.ValueOf(a)
	tb := reflect.ValueOf(b)

	if ta.Type().ConvertibleTo(tb.Type()) {
		return ta.Convert(tb.Type()).Interface()
	}
	return a
}

func compare(a, b interface{}) float64 {
	if reflect.DeepEqual(a, b) {
		return 0
	}

	if na, ok := castType(a, 0.0).(float64); ok {
		if nb, ok := castType(b, 0.0).(float64); ok {
			return na - nb
		}
	}

	sa := fmt.Sprintf("%#v", a)
	sb := fmt.Sprintf("%#v", b)

	return float64(strings.Compare(sa, sb))
}

func isNil(a interface{}) (yes bool) {
	if a == nil {
		return true
	}

	try(func() { yes = reflect.ValueOf(a).IsNil() })
	return
}

func (as Assertions) diff(a, b interface{}) string {
	if as.df != nil {
		return as.df(a, b)
	}
	return ""
}
