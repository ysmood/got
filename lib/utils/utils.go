package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func castType(x, y interface{}) interface{} {
	ta := reflect.ValueOf(x)
	tb := reflect.ValueOf(y)

	if (x == nil || y == nil) && (x != y) {
		return x
	}

	if ta.Type().ConvertibleTo(tb.Type()) {
		return ta.Convert(tb.Type()).Interface()
	}
	return x
}

// Compare returns the float value of x minus y
func Compare(x, y interface{}) float64 {
	if reflect.DeepEqual(x, y) {
		return 0
	}

	if na, ok := castType(x, 0.0).(float64); ok {
		if nb, ok := castType(y, 0.0).(float64); ok {
			return na - nb
		}
	}

	sa := fmt.Sprintf("%#v", x)
	sb := fmt.Sprintf("%#v", y)

	return float64(strings.Compare(sa, sb))
}
