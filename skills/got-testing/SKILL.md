---
name: got-testing
description: Write Go tests with the `got` test framework (github.com/ysmood/got) — fluent assertions, setup helpers, helper-func-driven cases, test suites, HTTP serving/requests, value snapshots, mocking interfaces, and customizing assertion output. Use when adding or editing `*_test.go` files in any project that imports `github.com/ysmood/got`.
---

# Testing with `got`

`got` is a fluent Go test framework. You get a context `g` from `got.T(t)` (or a
custom `setup(t)` helper) that embeds the standard `*testing.T`, an assertions
API, and testing utilities.

To discover the full API at any time:

```shell
go doc github.com/ysmood/got.Assertions   # all assertion methods
go doc github.com/ysmood/got.Utils        # all helper utilities
go doc github.com/ysmood/got.G            # the test context
```

## 1. Simple assertion

Use `got` as a lightweight assertion lib inside a plain `testing.T` function.
Tests live in package `<pkg>_test`.

```go
package example_test

import (
	"testing"

	"github.com/ysmood/got"
)

func TestAssertion(t *testing.T) {
	got.T(t).Eq(1+1, 2)
}
```

`got.T(t)` is shorthand for `got.New(t)`. Both return a `got.G` context.

## 2. Fluent chaining

Every assertion is chainable. Common modifiers:

- `.Desc("...")` — attach a message printed on failure.
- `.Must()` — abort the test immediately (`FailNow`) instead of continuing.

```go
g := got.T(t)
g.Desc("1 + 2 must be 3").Must().Eq(1+2, 3)
```

### Assertion cheat sheet

| Method | Asserts |
|---|---|
| `Eq(x, y)` / `Neq(x, y)` | deep/smart equal (auto numeric coercion) / not equal |
| `Equal(x, y)` | strict `==` equality (no coercion) |
| `Gt`/`Gte`/`Lt`/`Lte(x, y)` | ordered comparison |
| `InDelta(x, y, delta)` | numbers within `delta` |
| `True(b)` / `False(b)` | boolean |
| `Nil(...)` / `NotNil(...)` | last arg is nil / not nil |
| `Zero(x)` / `NotZero(x)` | zero value of its type |
| `E(...)` | last arg is a nil error (great for `g.E(os.Open(...))`) |
| `Err(...)` | last arg is a non-nil error |
| `Is(x, y)` | `errors.Is` / same type |
| `Has(container, item)` | string/slice/map contains item |
| `Len(list, n)` | length equals `n` |
| `Regex(pattern, str)` | string matches regexp |
| `Panic(fn)` | `fn` panics; returns the recovered value |
| `Count(n)` | returns a fn; test fails unless it's called exactly `n` times |

```go
g.E(os.Open("go.mod"))                 // no error
g.Has([]int{1, 2, 3}, 2)               // contains
g.Len("abc", 3)                        // length
val := g.Panic(func() { panic("boom") })
g.Eq(val, "boom")
```

## 3. The `setup` helper pattern

Centralize per-test config (parallelism, timeouts, custom fields) in one `setup`
function, then start every test with `g := setup(t)`.

```go
func init() {
	// Set a default "go test" flag if not already provided.
	got.DefaultFlags("timeout=10s")
}

// G is your custom context — embed got.G and add fields you want per test.
type G struct {
	got.G
	now string
}

var setup = func(t *testing.T) G {
	g := got.T(t)

	g.Parallel()                  // run tests concurrently
	g.PanicAfter(time.Second)     // per-test timeout
	g.Cleanup(func() {})          // runs after the test, always

	return G{g, time.Now().Format(time.DateTime)}
}

func TestSetup(t *testing.T) {
	g := setup(t)
	g.Gt(g.now, "2023-01-02 15:04:05")
}
```

## 4. Helper-func-driven tests (prefer over table-driven)

For repeated cases, define a local `check` closure and call it once per case.
Each call is a distinct line, so a failure points at the exact call site and you
can set a breakpoint or step through a single case. Mark the closure with
`g.Helper()` so failures report the caller's line, not the closure's.

```go
func TestHelperFuncDriven(t *testing.T) {
	g := setup(t)

	check := func(a, b, want int) {
		g.Helper()
		g.Eq(a+b, want)
	}

	check(1, 2, 3)
	check(2, 3, 5)
}
```

Avoid table-driven tests (`[]struct{...}` + `for` loop): the shared loop body
hides which row failed, breakpoints fire on every iteration, and the indirection
hurts readability. Only reach for them when cases truly number in the dozens.

## 5. Test suites with `got.Each` (avoid if possible)

Prefer plain `TestXxx(t *testing.T)` functions with `setup(t)` — they're easier
to read, run individually, and debug. Reach for suites only when several cases
genuinely need to share the same struct fields/lifecycle.

Each exported method on a struct that embeds `got.G` becomes a test case.

```go
func TestSuite(t *testing.T) {
	got.Each(t, SumSuite{})
}

type SumSuite struct{ got.G }

func (g SumSuite) Sum() { g.Eq(1+1, 2) }
```

Pass a function to `got.Each` to initialize the context per case:

```go
got.Each(t, func(t *testing.T) AdvancedSuite {
	g := got.New(t)
	g.Parallel()
	g.PanicAfter(time.Second)
	return AdvancedSuite{g, 1, 2}
})

type AdvancedSuite struct {
	got.G
	a, b int
}

func (g AdvancedSuite) Sum() { g.Eq(g.a+g.b, 3) }
```

## 6. Utilities (`got.Utils`)

The context exposes handy helpers. Highlights:

- **HTTP server**: `s := g.Serve(); s.Mux.HandleFunc("/", handler)` then
  `s.URL("?a=1")`. The server is torn down automatically.
- **HTTP client**: `g.Req("", url).Bytes().String()` / `.JSON()` /
  `.Unmarshal(&v)`. Method defaults to GET when empty.
- **Goroutines**: `g.Go(fn)` waits for the goroutine before the test ends (a
  bare `go fn()` would be lost).
- **Filesystem & misc**: `g.WriteFile`, `g.Read`, `g.Open`, `g.PathExists`,
  `g.Chdir`, `g.Setenv`, `g.RandStr`, `g.RandInt`, `g.JSON`, `g.ToJSONString`.

```go
func TestUtils(t *testing.T) {
	g := setup(t)

	s := g.Serve()
	s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	val := g.Req("", s.URL("/")).Bytes().String()
	g.Eq(val, "pong")
}
```

## 7. Snapshot assertions

`g.Snapshot(name, value)` records the value on first run under
`.got/snapshots/<TestName>/<name>.gop` and fails on future mismatches. Delete
the file to re-record.

```go
g.Snapshot("the map value", map[int]string{1: "1", 2: "2"})
```

## 8. Mocking interfaces (`lib/mock`)

Embed `mock.Mock` in a struct, route each method through `mock.Proxy`, then stub
behavior.

```go
import "github.com/ysmood/got/lib/mock"

type mockWriter struct{ mock.Mock }

func (m *mockWriter) Write(b []byte) (int, error) {
	return mock.Proxy(m, m.Write)(b)
}

func TestMocking(t *testing.T) {
	g := setup(t)
	m := &mockWriter{}

	// Stub with a custom function.
	mock.Stub(m, m.Write, func(b []byte) (int, error) {
		g.Eq(string(b), "3")
		return 0, nil
	})

	// Or define input/output expectations.
	mock.On(m, m.Write).When([]byte("3")).Return(1, nil)

	_, _ = m.Write([]byte("3"))

	// Inspect call history (use with Snapshot).
	g.Len(m.Calls(m.Write), 1)
}
```

Use `m.Fallback(realImpl)` to delegate every non-stubbed method to a real
implementation — useful for interfaces with many methods.

## 9. Customizing assertion output

Override `g.Assertions.ErrorHandler` to change failure messages.

```go
import (
	"github.com/ysmood/gop"
	"github.com/ysmood/got/lib/diff"
)

dh := got.NewDefaultAssertionError(10, gop.ThemeDefault, diff.ThemeDefault)
g.Assertions.ErrorHandler = got.AssertionErrorReport(func(c *got.AssertionCtx) string {
	if c.Type == got.AssertionEq {
		return fmt.Sprintf("%v != %v", c.Details[0], c.Details[1])
	}
	return dh.Report(c)
})
```

## Running tests

```shell
go test ./...                                          # run everything
go test -run TestAssertion ./...                       # run one test by name
go test -race -coverprofile=coverage.out ./...         # with race + coverage
go run github.com/ysmood/got/cmd/check-cov@latest      # enforce coverage (100% by default)
```

## Conventions

- Put tests in package `<pkg>_test` and start each test with `g := setup(t)`
  (or `got.T(t)` for trivial cases).
- Prefer `got`'s fluent assertions over `if x != y { t.Fatal(...) }`.
- Prefer helper-func-driven cases over table-driven loops (section 4).
- Prefer plain test functions over suite style (`got.Each`) when possible
  (section 5).
- Run `go doc github.com/ysmood/got.Assertions` (or `.Utils` / `.G`) to explore
  the full API.
