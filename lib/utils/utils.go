package utils

import (
	"reflect"
	"strings"
	"time"

	"github.com/ysmood/got/lib/gop"
)

var float64Type = reflect.TypeOf(0.0)

// SmartCompare returns the float value of x minus y
func SmartCompare(x, y interface{}) float64 {
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

// MethodType of target method
func MethodType(target interface{}, method string) reflect.Type {
	targetVal := reflect.ValueOf(target)
	return targetVal.MethodByName(method).Type()
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
