package got

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

// Testable interface. Usually, you use *testing.T as it.
type Testable interface {
	Name() string                            // same as testing.common.Name
	Skipped() bool                           // same as testing.common.Skipped
	Failed() bool                            // same as testing.common.Failed
	Cleanup(func())                          // same as testing.common.Cleanup
	FailNow()                                // same as testing.common.FailNow
	Fail()                                   // same as testing.common.Fail
	Helper()                                 // same as testing.common.Helper
	Logf(format string, args ...interface{}) // same as testing.common.Logf
	SkipNow()                                // same as testing.common.Skip
}

// G is the helper context, it hold some useful helpers to write tests
type G struct {
	Testable
	Assertions
	Utils
}

// Options for Assertion
type Options struct {
	// Dump a value to human readable string
	Dump func(interface{}) string

	// Format keywords in the assertion message.
	// Such as color it for CLI output.
	Keyword func(string) string

	// Diff function for Assertions.Eq
	Diff func(a, b interface{}) string
}

var floatType = reflect.TypeOf(0.0)

// Defaults for Options
func Defaults() Options {
	return Options{
		func(v interface{}) (s string) {
			if v == nil {
				return "nil"
			}

			json := func() {
				buf := bytes.NewBuffer(nil)
				enc := json.NewEncoder(buf)
				enc.SetEscapeHTML(false)
				if enc.Encode(v) == nil {
					b, _ := json.Marshal(v)
					s = string(b)
				}
			}

			switch v.(type) {
			case string:
				json()
			case int:
				json()
			case bool:
				json()
			default:
				t := reflect.TypeOf(v)
				if t.ConvertibleTo(floatType) {
					s = fmt.Sprintf("%s(%v)", t, v)
				} else {
					s = fmt.Sprintf("%#v", v)
				}
			}

			return s
		},
		func(s string) string {
			return "⦗" + s + "⦘"
		},
		nil,
	}
}

// New G
func New(t Testable) G {
	return NewWith(t, Defaults())
}

// NewWith G with options
func NewWith(t Testable, opts Options) G {
	return G{t, Assertions{t, opts.Dump, opts.Keyword, opts.Diff}, Utils{t, os.Exit}}
}
