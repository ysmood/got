// Package utils ...
package utils

import (
	"reflect"
	"strings"
	"time"

	"github.com/ysmood/gop"
)

var float64Type = reflect.TypeOf(0.0)

// SmartCompare returns the float value of x minus y.
// If x and y are numerical types, the result will be the subtraction between them, such as x is int(1), y is float64(1.2),
// the result will be -0.2 . time.Time is also a numerical value.
// If x or y are not numerical types, both of them will be converted to string format of its value type, the result will be
// the strings.Compare result between them, such as x is int(1), y is "a", the result will be 1 .
func SmartCompare(x, y interface{}) float64 {
	_, xNil := IsNil(x)
	_, yNil := IsNil(y)
	if xNil && yNil {
		return 0
	}

	if reflect.DeepEqual(x, y) {
		return 0
	}

	if x != nil && y != nil {
		xVal := reflect.Indirect(reflect.ValueOf(x))
		yVal := reflect.Indirect(reflect.ValueOf(y))

		if xVal.Type().ConvertibleTo(float64Type) && yVal.Type().ConvertibleTo(float64Type) {
			return xVal.Convert(float64Type).Float() - yVal.Convert(float64Type).Float()
		}

		if xt, ok := xVal.Interface().(time.Time); ok {
			if yt, ok := yVal.Interface().(time.Time); ok {
				return float64(xt.Sub(yt))
			}
		}
	}

	return Compare(x, y)
}

// Compare returns the float value of x minus y
func Compare(x, y interface{}) float64 {
	return float64(strings.Compare(gop.Plain(x), gop.Plain(y)))
}

// ToInterfaces convertor
func ToInterfaces(vs []reflect.Value) []interface{} {
	out := []interface{}{}
	for _, v := range vs {
		out = append(out, v.Interface())
	}
	return out
}

// ToValues convertor
func ToValues(vs []interface{}) []reflect.Value {
	out := []reflect.Value{}
	for _, v := range vs {
		out = append(out, reflect.ValueOf(v))
	}
	return out
}

// IsNil returns true, true if the value is nilable and is nil.
func IsNil(x interface{}) (bool, bool) {
	if x == nil {
		return true, true
	}

	val := reflect.ValueOf(x)
	k := val.Kind()
	nilable := k == reflect.Chan ||
		k == reflect.Func ||
		k == reflect.Interface ||
		k == reflect.Map ||
		k == reflect.Ptr ||
		k == reflect.Slice

	if nilable {
		return true, val.IsNil()
	}

	return false, false
}
