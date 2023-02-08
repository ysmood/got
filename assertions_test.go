package got_test

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ysmood/gop"
	"github.com/ysmood/got"
)

func TestAssertion(t *testing.T) {
	as := setup(t)

	as.Eq(1, 1)
	as.Eq(1.0, 1)
	as.Eq([]int{1, 3}, []int{1, 3})
	as.Eq(map[int]int{1: 2, 3: 4}, map[int]int{3: 4, 1: 2})
	as.Eq(nil, nil)
	fn := func() {}
	as.Eq(map[int]interface{}{1: fn, 2: nil}, map[int]interface{}{2: nil, 1: fn})

	as.Neq(1.1, 1)
	as.Neq([]int{1, 2}, []int{2, 1})
	as.Neq("true", true)
	as.Neq(errors.New("a"), errors.New("b"))

	as.Equal(1, 1)
	arr := []int{1, 2}
	as.Equal(arr, arr)
	as.Equal(fn, fn)

	as.Lt(time.Millisecond, time.Second)
	as.Lte(1, 1)

	as.Gt(2, 1.5)
	as.Gte(2, 2.0)

	now := time.Now()
	as.Eq(now, now)
	as.Lt(now, now.Add(time.Second))
	as.Gt(now.Add(time.Second), now)

	as.InDelta(1.1, 1.2, 0.2)

	as.True(true)
	as.False(false)

	as.Nil(nil)
	as.Nil((*int)(nil))
	as.Nil(os.Stat("go.mod"))
	as.NotNil([]int{})

	as.Zero("")
	as.Zero(0)
	as.Zero(time.Time{})
	as.NotZero(1)
	as.NotZero("ok")
	as.NotZero(time.Now())

	as.Regex(`\d\d`, "10")
	as.Has(`test`, 'e')
	as.Has(`test`, "es")
	as.Has(`test`, []byte("es"))
	as.Has([]byte(`test`), "es")
	as.Has([]byte(`test`), []byte("es"))
	as.Has([]int{1, 2, 3}, 2)
	as.Has([3]int{1, 2, 3}, 2)
	as.Has(map[int]int{1: 4, 2: 5, 3: 6}, 5)

	as.Len([]int{1, 2}, 2)

	as.Err(1, 2, errors.New("err"))
	as.Panic(func() { panic(1) })

	as.Is(1, 2)
	err := errors.New("err")
	as.Is(err, err)
	as.Is(fmt.Errorf("%w", err), err)
	as.Is(nil, nil)

	as.Must().Eq(1, 1)

	count := as.Count(2)
	count()
	count()
}

func TestAssertionErr(t *testing.T) {
	m := &mock{t: t}
	as := got.New(m)
	as.Assertions.ErrorHandler = got.NewDefaultAssertionError(gop.ThemeNone, nil)

	type data struct {
		A int
		S string
	}

	as.Desc("not %s", "equal").Eq(1, 2.0)
	m.check("not equal\n1 ⦗not ==⦘ 2.0")

	as.Eq(data{1, "a"}, data{1, "b"})
	m.check(`
got_test.data{
    A: 1,
    S: "a",
}

⦗not ==⦘

got_test.data{
    A: 1,
    S: "b",
}`)

	as.Eq(true, "a&")
	m.check(`true ⦗not ==⦘ "a&"`)

	as.Eq(nil, "ok")
	m.check(`nil ⦗not ==⦘ "ok"`)

	as.Eq(1, nil)
	m.check(`1 ⦗not ==⦘ nil`)

	as.Equal(1, 1.0)
	m.check("1 ⦗not ==⦘ 1.0")
	as.Equal([]int{1}, []int{2})
	m.check(`
[]int/* len=1 cap=1 */{
    1,
}

⦗not ==⦘

[]int/* len=1 cap=1 */{
    2,
}`)

	as.Neq(1, 1)
	m.check("1 ⦗==⦘ 1")
	as.Neq(1.0, 1)
	m.check("1.0 ⦗==⦘ 1 ⦗when converted to the same type⦘ ")

	as.Lt(1, 1)
	m.check("1 ⦗not <⦘ 1")
	as.Lte(2, 1)

	m.check("2 ⦗not ≤⦘ 1")
	as.Gt(1, 1)
	m.check("1 ⦗not >⦘ 1")
	as.Gte(1, 2)
	m.check("1 ⦗not ≥⦘ 2")

	as.InDelta(10, 20, 3)
	m.check(" ⦗delta between⦘ 10 ⦗and⦘ 20 ⦗not ≤⦘ 3.0")

	as.True(false)
	m.check(" ⦗should be⦘ true")
	as.False(true)
	m.check(" ⦗should be⦘ false")

	as.Nil(1)
	m.check(" ⦗last argument⦘ 1 ⦗should be⦘ nil")
	as.Nil()
	m.check(" ⦗no arguments received⦘ ")
	as.NotNil(nil)
	m.check(" ⦗last argument shouldn't be⦘ nil")
	as.NotNil((*int)(nil))
	m.check(" ⦗last argument⦘ (*int)(nil) ⦗shouldn't be⦘ nil")
	as.NotNil()
	m.check(" ⦗no arguments received⦘ ")
	as.NotNil(1)
	m.check(" ⦗last argument⦘ 1 ⦗is not nilable⦘ ")

	as.Zero(1)
	m.check("1 ⦗should be zero value for its type⦘ ")
	as.NotZero(0)
	m.check("0 ⦗shouldn't be zero value for its type⦘ ")

	as.Regex(`\d\d`, "aaa")
	m.check(`"\\d\\d" ⦗should match⦘ "aaa"`)
	as.Has(`test`, "x")
	m.check(`"test" ⦗should has⦘ "x"`)

	as.Len([]int{1, 2}, 3)
	m.check(" ⦗expect len⦘ 2 ⦗to be⦘ 3")

	as.Err(nil)
	m.check(" ⦗last value⦘ nil ⦗should be <error>⦘ ")
	as.Panic(func() {})
	m.check(" ⦗should panic⦘ ")
	as.Err()
	m.check(" ⦗no arguments received⦘ ")
	as.Err(1)
	m.check(" ⦗last value⦘ 1 ⦗should be <error>⦘ ")

	func() {
		defer func() {
			_ = recover()
		}()
		as.E(1, errors.New("E"))
	}()
	m.check(`
⦗last argument⦘

&errors.errorString{
    s: "E",
}

⦗should be⦘

nil`)

	as.Is(1, 2.2)
	m.check("1 ⦗should be kind of⦘ 2.2")
	as.Is(errors.New("a"), errors.New("b"))
	m.check(`
&errors.errorString{
    s: "a",
}

⦗should in chain of⦘

&errors.errorString{
    s: "b",
}`)
	as.Is(nil, errors.New("a"))
	m.check(`
nil

⦗should be kind of⦘

&errors.errorString{
    s: "a",
}`)
	as.Is(errors.New("a"), nil)
	m.check(`
&errors.errorString{
    s: "a",
}

⦗should be kind of⦘

nil`)

	{
		count := as.Count(2)
		count()
		m.cleanup()
		m.check(` ⦗should count⦘ 2 ⦗times, but got⦘ 1`)

		count = as.Count(1)
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			count()
			wg.Done()
		}()
		go func() {
			count()
			wg.Done()
		}()
		wg.Wait()
		m.cleanup()
		m.check(` ⦗should count⦘ 1 ⦗times, but got⦘ 2`)
	}
}

func TestAssertionColor(t *testing.T) {
	m := &mock{t: t}

	g := got.New(m)
	g.Eq([]int{1, 2}, []int{1, 3})
	m.checkWithStyle(true, `
<36>[]int<39><37>/* len=2 cap=2 */<39>{
    <32>1<39>,
    <32>2<39>,
}

<31><4>⦗not ==⦘<24><39>

<36>[]int<39><37>/* len=2 cap=2 */<39>{
    <32>1<39>,
    <32>3<39>,
}

<45><30>@@ diff chunk @@<39><49>
2 2       1,
<31>3   -<39>     <31>2<39>,
<32>  3 +<39>     <32>3<39>,
4 4   }

`)

	g.Eq("abc", "axc")
	m.checkWithStyle(true, `"a<31>b<39>c" <31><4>⦗not ==⦘<24><39> "a<32>x<39>c"`)

	g.Eq(3, "a")
	m.checkWithStyle(true, `<31>3<39> <31><4>⦗not ==⦘<24><39> <32>"a"<39>`)
}

func TestCustomAssertionError(t *testing.T) {
	m := &mock{t: t}

	g := got.New(m)
	g.ErrorHandler = got.AssertionErrorReport(func(c *got.AssertionCtx) string {
		if c.Type == got.AssertionEq {
			return "custom eq"
		}
		return ""
	})
	g.Eq(1, 2)
	m.check("custom eq")
}
