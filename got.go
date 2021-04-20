package got

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
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
					s = buf.String()[:buf.Len()-1]
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
			return " ⦗" + s + "⦘ "
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
	return G{
		t,
		Assertions{t, false, "", opts.Dump, opts.Keyword, opts.Diff},
		Utils{t},
	}
}

// DefaultFlags will set the "go test" flag if not yet presented.
// It must be executed in the init() function.
// Such as the timeout:
//     DefaultFlags("timeout=10s")
func DefaultFlags(flags ...string) {
	// remove default timeout from "go test"
	filtered := []string{}
	for _, arg := range os.Args {
		if arg != "-test.timeout=10m0s" {
			filtered = append(filtered, arg)
		}
	}
	os.Args = filtered

	list := map[string]struct{}{}
	reg := regexp.MustCompile(`^-test\.(\w+)`)
	for _, arg := range os.Args {
		ms := reg.FindStringSubmatch(arg)
		if ms != nil {
			list[ms[1]] = struct{}{}
		}
	}

	for _, flag := range flags {
		if _, has := list[strings.Split(flag, "=")[0]]; !has {
			os.Args = append(os.Args, "-test."+flag)
		}
	}
}

// Parallel config of "go test -parallel"
func Parallel() (n int) {
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "test.parallel" {
			v := reflect.ValueOf(f.Value).Elem().Convert(reflect.TypeOf(n))
			n = v.Interface().(int)
		}
	})
	return
}
