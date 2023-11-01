// Package mock provides a simple way to stub struct methods.
package mock

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/ysmood/got/lib/utils"
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
}

// Fallback the methods that are not stubbed to fb.
func (m *Mock) Fallback(fb interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.fallback = reflect.ValueOf(fb)
}

// Stub the method with stub
func Stub[M any](mock Fallbackable, method M, stub M) {
	panicIfNotFunc(method)

	m := toMock(mock)

	m.lock.Lock()
	defer m.lock.Unlock()

	if m.stubs == nil {
		m.stubs = map[string]interface{}{}
	}

	m.stubs[fnName(method)] = stub
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

	return m.fallback.MethodByName(name).Interface().(M)
}

// StubOn utils
type StubOn struct {
	when []*StubWhen
}

// StubWhen utils
type StubWhen struct {
	lock  *sync.Mutex
	on    *StubOn
	in    []interface{}
	ret   *StubReturn
	count int // how many times this stub has been matched
}

// StubReturn utils
type StubReturn struct {
	on    *StubOn
	out   []reflect.Value
	times *StubTimes
}

// StubTimes utils
type StubTimes struct {
	count int
}

// On helper to stub methods to conditionally return values.
func On[M any](mock Fallbackable, method M) *StubOn {
	panicIfNotFunc(method)

	m := toMock(mock)

	s := &StubOn{
		when: []*StubWhen{},
	}

	eq := func(in, arg []interface{}) bool {
		for i := 0; i < len(in); i++ {
			if in[i] != Any && utils.Compare(in[i], arg[i]) != 0 {
				return false
			}
		}
		return true
	}

	t := reflect.TypeOf(method)

	fn := reflect.MakeFunc(t, func(args []reflect.Value) []reflect.Value {
		argsIt := utils.ToInterfaces(args)
		for _, when := range s.when {
			if eq(when.in, argsIt) {
				when.lock.Lock()

				when.ret.times.count--
				if when.ret.times.count == 0 {
					m.Stop(method)
				}

				when.count++

				when.lock.Unlock()

				return toReturnValues(t, when.ret.out)
			}
		}
		panic(fmt.Sprintf("No mock.StubOn.When matches: %#v", argsIt))
	})

	Stub(m, method, fn.Interface().(M))

	return s
}

// Any input
var Any = struct{}{}

// When input args of stubbed method matches in
func (s *StubOn) When(in ...interface{}) *StubWhen {
	w := &StubWhen{lock: &sync.Mutex{}, on: s, in: in}
	s.when = append(s.when, w)
	return w
}

// Return the out as the return values of stubbed method
func (s *StubWhen) Return(out ...interface{}) *StubReturn {
	r := &StubReturn{on: s.on, out: utils.ToValues(out)}
	r.Times(0)
	s.ret = r
	return r
}

// Count returns how many times this condition has been matched
func (s *StubWhen) Count() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.count
}

// Times specifies how how many stubs before stop, if n <= 0 it will never stop.
func (s *StubReturn) Times(n int) *StubOn {
	t := &StubTimes{count: n}
	s.times = t
	return s.on
}

// Once specifies stubs only once before stop
func (s *StubReturn) Once() *StubOn {
	return s.Times(1)
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
