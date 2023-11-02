package mock

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/ysmood/got/lib/utils"
)

// Stub the method with stub
func Stub[M any](mock Fallbackable, method M, stub M) {
	panicIfNotFunc(method)

	m := toMock(mock)

	m.lock.Lock()
	defer m.lock.Unlock()

	if m.stubs == nil {
		m.stubs = map[string]interface{}{}
	}

	name := fnName(method)

	m.stubs[name] = m.spy(name, stub)
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
