package utils_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ysmood/got/lib/utils"
)

func TestSmartCompare(t *testing.T) {
	now := time.Now()

	circular := map[int]interface{}{}
	circular[0] = circular

	fn := func() {}
	fn2 := func() {}
	ch := make(chan int, 1)
	ch2 := make(chan int, 1)

	testCases := []struct {
		x interface{}
		y interface{}
		s interface{}
	}{
		{1, 1, 0.0},
		{1, 3.0, -2.0},
		{1, "a", 1.0},
		{"b", "a", 1.0},
		{1, nil, -1.0},
		{fn, fn, 0.0},
		{fn, fn2, -1.0},
		{ch, ch, 0.0},
		{ch, ch2, -1.0},
		{now.Add(time.Second), now, float64(time.Second)},
		{circular, circular, 0.0},
		{circular, 0, 1.0},
		{map[int]interface{}{1: 1.0}, map[int]interface{}{1: 1}, 1.0},
	}
	for i, c := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			s := utils.SmartCompare(c.x, c.y)
			if s != c.s {
				t.Fail()
				t.Log("expect s to be", c.s, "but got", s)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	if utils.Compare(1, 1.0) == 0 {
		t.Fail()
	}
}

func TestOthers(t *testing.T) {
	fn := reflect.MakeFunc(utils.MethodType(t, "Name"), func(args []reflect.Value) (results []reflect.Value) {
		return []reflect.Value{reflect.ValueOf("test")}
	}).Interface()

	if fn.(func() string)() != "test" {
		t.Error("fail")
	}

	vs := utils.ToValues([]interface{}{1})

	if utils.ToInterfaces(vs)[0] != 1 {
		t.Error("fail")
	}
}
