package got

import (
	"flag"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/ysmood/got/lib/diff"
	"github.com/ysmood/got/lib/gop"
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

// G is the helper context, it provides some handy helpers for testing
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

func defaults() Options {
	if _, has := os.LookupEnv("NO_COLOR"); has {
		return NoColor()
	}
	return Defaults()
}

// Defaults for Options
func Defaults() Options {
	return Options{
		func(v interface{}) string {
			return gop.F(v)
		},
		func(s string) string {
			return gop.ColorStr(gop.Red, " ⦗"+s+"⦘ ")
		},
		func(a, b interface{}) string {
			df := diff.Diff(gop.Plain(a), gop.Plain(b))
			return "\n\n" + df + "\n"
		},
	}
}

// NoColor defaults for Options
func NoColor() Options {
	return Options{
		func(v interface{}) string {
			return gop.Plain(v)
		},
		func(s string) string {
			return " ⦗" + s + "⦘ "
		},
		func(a, b interface{}) string {
			df := diff.Diff(gop.Plain(a), gop.Plain(b))
			return "\n\n" + df + "\n"
		},
	}
}

// NoDiff returns a clone with diff output disabled.
func (opts Options) NoDiff() Options {
	opts.Diff = nil
	return opts
}

// Setup returns a helper to init G instance.
// It will respect https://no-color.org/
func Setup(init func(g G)) func(t Testable) G {
	return SetupWith(defaults(), init)
}

// SetupWith options
func SetupWith(opts Options, init func(g G)) func(t Testable) G {
	return func(t Testable) G {
		g := NewWith(t, opts)
		if init != nil {
			init(g)
		}
		return g
	}
}

// T is the shortcut for New
func T(t Testable) G {
	return New(t)
}

// New G instance.
// It will respect https://no-color.org/
func New(t Testable) G {
	return NewWith(t, defaults())
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
