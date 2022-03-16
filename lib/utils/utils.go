package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

var float64Type = reflect.TypeOf(0.0)

// Compare returns the float value of x minus y
func Compare(x, y interface{}) float64 {
	if reflect.DeepEqual(x, y) {
		return 0
	}

	if x != nil && y != nil {
		xVal := reflect.Indirect(reflect.ValueOf(x))
		yVal := reflect.Indirect(reflect.ValueOf(y))

		if xVal.CanConvert(float64Type) && yVal.CanConvert(float64Type) {
			return xVal.Convert(float64Type).Float() - yVal.Convert(float64Type).Float()
		}

		if xt, ok := xVal.Interface().(time.Time); ok {
			if yt, ok := yVal.Interface().(time.Time); ok {
				return float64(xt.Sub(yt))
			}
		}
	}

	sa := fmt.Sprintf("%#v", x)
	sb := fmt.Sprintf("%#v", y)

	return float64(strings.Compare(sa, sb))
}
