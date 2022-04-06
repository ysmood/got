package gop

import (
	"reflect"
	"unsafe"
)

// GetPrivateField via field index
// TODO: we can use a LRU cache for the copy of the values, but it might be trivial for just testing.
func GetPrivateField(v reflect.Value, i int) reflect.Value {
	if v.Kind() != reflect.Struct {
		panic("expect v to be a struct")
	}

	copied := reflect.New(v.Type()).Elem()
	copied.Set(v)
	f := copied.Field(i)
	return reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
}

// GetPrivateFieldByName is similar with GetPrivateField
func GetPrivateFieldByName(v reflect.Value, name string) reflect.Value {
	if v.Kind() != reflect.Struct {
		panic("expect v to be a struct")
	}

	copied := reflect.New(v.Type()).Elem()
	copied.Set(v)
	f := copied.FieldByName(name)
	return reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
}
