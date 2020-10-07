package got

import (
	"fmt"
	"reflect"
)

// Testable interface
type Testable interface {
	Helper()
	Fail()
	FailNow()
	Cleanup(func())
	Logf(format string, args ...interface{})
}

// G is the helper context
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
			return fmt.Sprintf("%v (%v)", v, reflect.TypeOf(v))
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
