package got

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

// Assertion is the assertion context
type Assertion struct {
	Testable
}

// New assertion helper
func New(t Testable) Assertion {
	return Assertion{t}
}

// Eq a == b
func (as Assertion) Eq(a, b interface{}) {
	as.Helper()
	if compare(a, b) != 0 {
		as.err("%s == %s", pp(a), pp(b))
	}
}

// Neq a != b
func (as Assertion) Neq(a, b interface{}) {
	as.Helper()
	if compare(a, b) == 0 {
		as.err("%s != %s", pp(a), pp(b))
	}
}

// Gt a > b
func (as Assertion) Gt(a, b interface{}) {
	as.Helper()
	if compare(a, b) <= 0 {
		as.err("%s > %s", pp(a), pp(b))
	}
}

// Gte a >= b
func (as Assertion) Gte(a, b interface{}) {
	as.Helper()
	if compare(a, b) < 0 {
		as.err("%s >= %s", pp(a), pp(b))
	}
}

// Lt a < b
func (as Assertion) Lt(a, b interface{}) {
	as.Helper()
	if compare(a, b) >= 0 {
		as.err("%s < %s", pp(a), pp(b))
	}
}

// Lte a <= b
func (as Assertion) Lte(a, b interface{}) {
	as.Helper()
	if compare(a, b) > 0 {
		as.err("%s <= %s", pp(a), pp(b))
	}
}

// True a == true
func (as Assertion) True(a bool) {
	as.Helper()
	if !a {
		as.err("should be true")
	}
}

// False a == false
func (as Assertion) False(a bool) {
	as.Helper()
	if a {
		as.err("should be false")
	}
}

// Nil args[-1] == nil
func (as Assertion) Nil(args ...interface{}) {
	as.Helper()
	if len(args) == 0 {
		as.err("no args received")
		return
	}
	last := args[len(args)-1]
	if !isNil(last) {
		as.err("%s should be nil", pp(last))
	}
}

// NotNil args[-1] != nil
func (as Assertion) NotNil(args ...interface{}) {
	as.Helper()
	if len(args) == 0 {
		as.err("no args received")
		return
	}
	last := args[len(args)-1]
	if isNil(last) {
		as.err("%s shouldn't be nil", pp(last))
	}
}

// Regex matches str
func (as Assertion) Regex(pattern, str string) {
	as.Helper()
	if !regexp.MustCompile(pattern).MatchString(str) {
		as.err("%s <regex should match> %s", pattern, str)
	}
}

// Has str in container
func (as Assertion) Has(container, str string) {
	as.Helper()
	if !strings.Contains(container, str) {
		as.err("%s <should has> %s", container, str)
	}
}

// Len len(list) == l
func (as Assertion) Len(list interface{}, l int) {
	as.Helper()
	actual := reflect.ValueOf(list).Len()
	if actual != l {
		as.err("expect len %d to be %d", actual, l)
	}
}

// Err args[-1] is error and not nil
func (as Assertion) Err(args ...interface{}) {
	as.Helper()
	if len(args) == 0 {
		as.err("no args received")
		return
	}
	last := args[len(args)-1]
	if err, _ := last.(error); err == nil {
		as.err("%s should be error", pp(last))
	}
}

// Panic fn should panic
func (as Assertion) Panic(fn func()) {
	as.Helper()

	defer func() {
		as.Helper()

		val := recover()
		if val == nil {
			as.err("should panic")
		}
	}()

	fn()
}

// Is a a kind of b
func (as Assertion) Is(a, b interface{}) {
	as.Helper()

	if ae, ok := a.(error); ok {
		if be, ok := b.(error); ok {
			if ae == be {
				return
			}

			if errors.Is(ae, be) {
				return
			}
			as.err("%s <not in err chain> %s", pp(a), pp(b))
		}
	}

	at := reflect.TypeOf(a)
	bt := reflect.TypeOf(b)
	if at.Kind() != bt.Kind() {
		as.err("%s <not kind of> %s", pp(a), pp(b))
	}
}
