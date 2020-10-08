package got

import (
	"bytes"
	"encoding/json"
	"reflect"
)

func (u Assertions) err(format string, args ...interface{}) Result {
	u.Helper()
	u.Logf(format, args...)
	u.Fail()
	return Result{u, true}
}

func (hp Helpers) err(err error) {
	if err != nil {
		hp.Fatal(err)
	}
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

func isNil(a interface{}) (yes bool) {
	if a == nil {
		return true
	}

	try(func() { yes = reflect.ValueOf(a).IsNil() })
	return
}

func try(fn func()) {
	defer func() {
		_ = recover()
	}()
	fn()
}
