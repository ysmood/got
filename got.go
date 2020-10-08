package got

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

// Testable interface. Usually, you use *testing.T as it.
type Testable interface {
	Helper()                                 // same as testing.common.Helper
	Fail()                                   // same as testing.common.Fail
	FailNow()                                // same as testing.common.FailNow
	SkipNow()                                // same as testing.common.Skip
	Cleanup(func())                          // same as testing.common.Cleanup
	Logf(format string, args ...interface{}) // same as testing.common.Logf
}

// G is the helper context, it hold some useful helpers to write tests
type G struct {
	Testable

	d func(v interface{}) string // Options.Dump
	k func(string) string        // Options.Keyword
}

// Options for Assertion
type Options struct {
	// Dump a value to human readable string
	Dump func(interface{}) string

	// Format keywords in the assertion message.
	// Such as color it for CLI output.
	Keyword func(string) string
}

// Defaults for Options
func Defaults() Options {
	return Options{
		func(v interface{}) string {
			if v == nil {
				return "nil"
			}

			s := fmt.Sprintf("%v", v)

			json := func() {
				buf := bytes.NewBuffer(nil)
				enc := json.NewEncoder(buf)
				enc.SetEscapeHTML(false)
				if enc.Encode(v) == nil {
					b, _ := json.Marshal(v)
					s = string(b)
				}
			}

			t := ""
			switch v.(type) {
			case string:
				json()
			case int:
				json()
			case bool:
				json()
			default:
				t = fmt.Sprintf(" <%v>", reflect.TypeOf(v))
			}

			return s + t
		},
		func(s string) string {
			return "⦗" + s + "⦘"
		},
	}
}

// New assertion helper
func New(t Testable) G {
	return NewWith(t, Defaults())
}

// NewWith assertion helper with options
func NewWith(t Testable, opts Options) G {
	return G{t, opts.Dump, opts.Keyword}
}
