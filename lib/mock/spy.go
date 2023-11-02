package mock

import "reflect"

// Call record the input and output of a method call
type Call struct {
	Input  []any
	Return []any
}

// Calls returns all the calls of method
func (m *Mock) Calls(method any) []Call {
	panicIfNotFunc(method)

	m.lock.Lock()
	defer m.lock.Unlock()

	return m.calls[fnName(method)]
}

// Record all the input and output of a method
func (m *Mock) spy(name string, fn any) any {
	v := reflect.ValueOf(fn)
	t := v.Type()

	if m.calls == nil {
		m.calls = map[string][]Call{}
	}

	return reflect.MakeFunc(t, func(args []reflect.Value) []reflect.Value {
		ret := v.Call(args)

		m.lock.Lock()
		m.calls[name] = append(m.calls[name], Call{
			valToInterface(args),
			valToInterface(ret),
		})
		m.lock.Unlock()

		return ret
	}).Interface()
}

func valToInterface(list []reflect.Value) []any {
	ret := make([]any, len(list))
	for i, v := range list {
		ret[i] = v.Interface()
	}
	return ret
}
