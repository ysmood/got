package got

import (
	"reflect"
)

const prefixEachErr = "[got.Each]"

// Each runs each exported method Fn on type Ctx as a subtest of t.
// The iteratee can be a struct Ctx or:
//
//     iteratee(t Testable) (ctx Ctx)
//
// Each Fn will be called like:
//
//      ctx.Fn()
//
// If iteratee is Ctx, its Assertion and T fields will be set to New(t) and t for each test.
// Any Fn that has the same name with the embedded one will be ignored.
func Each(t Testable, iteratee interface{}) (count int) {
	t.Helper()

	itVal := normalizeIteratee(t, iteratee)

	ctxType := itVal.Type().Out(0)

	methods := filterMethods(ctxType)

	runVal := reflect.ValueOf(t).MethodByName("Run")
	cbType := runVal.Type().In(1)

	for _, m := range methods {
		checkFnType(t, m)

		// because the callback is in another goroutine, we create closures for each loop
		method := m

		runVal.Call([]reflect.Value{
			reflect.ValueOf(method.Name),
			reflect.MakeFunc(cbType, func(args []reflect.Value) []reflect.Value {
				res := itVal.Call(args)
				return method.Func.Call(res)
			}),
		})
	}
	return len(methods)
}

func normalizeIteratee(t Testable, iteratee interface{}) reflect.Value {
	t.Helper()

	if iteratee == nil {
		t.Logf(prefixEachErr + " iteratee shouldn't be nil")
		t.FailNow()
	}

	itVal := reflect.ValueOf(iteratee)
	itType := itVal.Type()
	fail := true

	switch itType.Kind() {
	case reflect.Func:
		if itType.NumIn() != 1 || itType.NumOut() != 1 {
			break
		}
		try(func() {
			_ = reflect.New(itType.In(0).Elem()).Interface().(Testable)
			fail = false
		})

	case reflect.Struct:
		fnType := reflect.FuncOf([]reflect.Type{reflect.TypeOf(t)}, []reflect.Type{itType}, false)
		structVal := itVal
		itVal = reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
			sub := args[0].Interface().(Testable)
			as := reflect.ValueOf(New(sub))

			c := reflect.New(itType).Elem()
			c.Set(structVal)
			try(func() { c.FieldByName("Assertion").Set(as) })
			try(func() { c.FieldByName("T").Set(args[0]) })

			return []reflect.Value{c}
		})
		fail = false
	}

	if fail {
		t.Logf(prefixEachErr+" iteratee <%v> should be a struct or <func(got.Testable) Ctx>", itType)
		t.FailNow()
	}
	return itVal
}

func checkFnType(t Testable, fn reflect.Method) {
	t.Helper()

	if fn.Type.NumIn() != 1 || fn.Type.NumOut() != 0 {
		t.Logf(prefixEachErr+" %s.%s shouldn't have arguments or return values", fn.Type.In(0).String(), fn.Name)
		t.FailNow()
	}
}

func filterMethods(typ reflect.Type) []reflect.Method {
	skip := map[string]struct{}{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous {
			for j := 0; j < field.Type.NumMethod(); j++ {
				skip[field.Type.Method(j).Name] = struct{}{}
			}
		}
	}

	methods := []reflect.Method{}
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		if _, has := skip[method.Name]; has {
			continue
		}
		methods = append(methods, method)
	}
	return methods
}
