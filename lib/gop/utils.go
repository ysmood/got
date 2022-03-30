package gop

import (
	"reflect"
	"regexp"
	"unsafe"
)

var regNewline = regexp.MustCompile(`\n`)

// GetPrivateField field value via field index
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
