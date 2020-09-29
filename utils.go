package got

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

// Testable interface
type Testable interface {
	Helper()
	Fail()
	FailNow()
	Logf(format string, args ...interface{})
}

func (as Assertion) err(format string, args ...interface{}) {
	as.Helper()
	as.Logf(format, args...)
	as.Fail()
}

// pretty print a value
func pp(v interface{}) string {
	return fmt.Sprintf("%v (%v)", v, reflect.TypeOf(v))
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

	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)

	return float64(bytes.Compare(ja, jb))
}

func isNil(a interface{}) bool {
	if a == nil {
		return true
	}

	defer func() { _ = recover() }()
	return reflect.ValueOf(a).IsNil()
}

func try(fn func()) {
	defer func() {
		_ = recover()
	}()
	fn()
}
