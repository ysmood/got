// Package mock provides a simple way to stub struct methods.
package mock

import (
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// Fallbackable interface
type Fallbackable interface {
	Fallback(any)
}

// Mock helper for interface stubbing
type Mock struct {
	lock sync.Mutex

	fallback reflect.Value

	stubs map[string]interface{}

	calls map[string][]Call
}

// Fallback the methods that are not stubbed to fb.
func (m *Mock) Fallback(fb interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.fallback = reflect.ValueOf(fb)
}

// Stop the stub
func (m *Mock) Stop(method any) {
	panicIfNotFunc(method)

	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.stubs, fnName(method))
}

// Proxy the input and output of method on mock for later stub.
func Proxy[M any](mock Fallbackable, method M) M {
	panicIfNotFunc(method)

	m := toMock(mock)

	m.lock.Lock()
	defer m.lock.Unlock()

	name := fnName(method)

	if fn, has := m.stubs[name]; has {
		return fn.(M)
	}

	if !m.fallback.IsValid() {
		panic("you should specify the mock.Mock.Fallback")
	}

	methodVal := m.fallback.MethodByName(name)
	if !methodVal.IsValid() {
		panic(m.fallback.Type().String() + " doesn't have method: " + name)
	}

	return m.spy(name, m.fallback.MethodByName(name).Interface()).(M)
}

func toMock(mock Fallbackable) *Mock {
	if m, ok := mock.(*Mock); ok {
		return m
	}

	return reflect.Indirect(reflect.ValueOf(mock)).FieldByName("Mock").Addr().Interface().(*Mock)
}

func fnName(fn interface{}) string {
	fv := reflect.ValueOf(fn)

	fi := runtime.FuncForPC(fv.Pointer())

	name := regexp.MustCompile(`^.+\.`).ReplaceAllString(fi.Name(), "")

	// remove the "-fm" suffix for struct methods
	name = strings.TrimSuffix(name, "-fm")

	return name
}

func panicIfNotFunc(fn any) {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		panic("the input should be a function")
	}
}

func toReturnValues(t reflect.Type, res []reflect.Value) []reflect.Value {
	out := []reflect.Value{}
	for i := 0; i < t.NumOut(); i++ {
		v := reflect.New(t.Out(i)).Elem()
		if res[i].IsValid() {
			v.Set(res[i])
		}
		out = append(out, v)
	}

	return out
}
