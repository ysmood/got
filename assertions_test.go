package got_test

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ysmood/got"
)

func TestAssertion(t *testing.T) {
	as := got.New(t)

	as.Eq(1, 1)
	as.Eq(1.0, 1)
	as.Eq([]int{1, 3}, []int{1, 3})
	as.Eq(map[int]int{1: 2, 3: 4}, map[int]int{3: 4, 1: 2})
	as.Eq(nil, nil)

	as.Neq(1.1, 1)
	as.Neq([]int{1, 2}, []int{2, 1})
	as.Neq("true", true)
	as.Neq(errors.New("a"), errors.New("b"))

	as.Equal(1, 1)

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
	as.Has(`test`, "es")

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
	as := got.NewWith(m, got.NoColor().NoDiff())

	type data struct {
		A int
		S string
	}

	as.Desc("not %s", "equal").Eq(1, 2.0)
	m.check("not equal\n1 ⦗not ==⦘ float64(2) ⦗even when converted to the same type⦘ ")

	as.Eq(data{1, "a"}, data{1, "b"})
	m.check(`got_test.data/* len=2 */{
    A: 1,
    S: "a",
} ⦗not ==⦘ got_test.data/* len=2 */{
    A: 1,
    S: "b",
}`)

	as.Eq(true, "a&")
	m.check(`true ⦗not ==⦘ "a&" ⦗even when converted to the same type⦘ `)

	as.Eq(nil, "ok")
	m.check(`nil ⦗not ==⦘ "ok" ⦗even when converted to the same type⦘ `)

	as.Eq(1, nil)
	m.check(`1 ⦗not ==⦘ nil ⦗even when converted to the same type⦘ `)

	as.Equal(1, 1.0)
	m.check("1 ⦗not ==⦘ float64(1)")

	as.Neq(1, 1)
	m.check("1 ⦗==⦘ 1")
	as.Neq(1.0, 1)
	m.check("float64(1) ⦗==⦘ 1 ⦗when converted to the same type⦘ ")

	as.Lt(1, 1)
	m.check("1 ⦗not <⦘ 1")
	as.Lte(2, 1)

	m.check("2 ⦗not ≤⦘ 1")
	as.Gt(1, 1)
	m.check("1 ⦗not >⦘ 1")
	as.Gte(1, 2)
	m.check("1 ⦗not ≥⦘ 2")

	as.InDelta(10, 20, 3)
	m.check("delta between 10 and 20 ⦗not ≤⦘ float64(3)")

	as.True(false)
	m.check(" ⦗should be⦘ true")
	as.False(true)
	m.check(" ⦗should be⦘ false")

	as.Nil(1)
	m.check(" ⦗last item in args⦘ 1 ⦗should be⦘ nil")
	as.Nil()
	m.check(" ⦗no args received⦘ ")
	as.NotNil(nil)
	m.check(" ⦗last value shouldn't be⦘ nil")
	as.NotNil((*int)(nil))
	m.check(" ⦗last item in args⦘ (*int)(nil) ⦗shouldn't be⦘ nil")
	as.NotNil()
	m.check(" ⦗no args received⦘ ")
	as.NotNil(1)
	m.check(" ⦗last item in args⦘ 1 ⦗is not nilable⦘ ")

	as.Zero(1)
	m.check("1 ⦗should be zero value for its type⦘ ")
	as.NotZero(0)
	m.check("0 ⦗should not be zero value for its type⦘ ")

	as.Regex(`\d\d`, "aaa")
	m.check(`\d\d ⦗should match⦘ aaa`)
	as.Has(`test`, "x")
	m.check("test ⦗should has⦘ x")

	as.Len([]int{1, 2}, 3)
	m.check(" ⦗expect len⦘ 2 ⦗to be⦘ 3")

	as.Err(nil)
	m.check(" ⦗last value⦘ nil ⦗should be <error>⦘ ")
	as.Panic(func() {})
	m.check(" ⦗should panic⦘ ")
	as.Err()
	m.check(" ⦗no args received⦘ ")
	as.Err(1)
	m.check(" ⦗last value⦘ 1 ⦗should be <error>⦘ ")

	func() {
		defer func() {
			_ = recover()
		}()
		as.E(1, errors.New("E"))
	}()
	m.check(` ⦗last item in args⦘ &errors.errorString{
    s: "E",
} ⦗should be⦘ nil`)

	as.Is(1, 2.2)
	m.check("1 ⦗should be kind of⦘ float64(2.2)")
	as.Is(errors.New("a"), errors.New("b"))
	m.check(`&errors.errorString{
    s: "a",
} ⦗should in chain of⦘ &errors.errorString{
    s: "b",
}`)
	as.Is(nil, errors.New("a"))
	m.check(`nil ⦗should be kind of⦘ &errors.errorString{
    s: "a",
}`)
	as.Is(errors.New("a"), nil)
	m.check(`&errors.errorString{
    s: "a",
} ⦗should be kind of⦘ nil`)

	opts := got.NoColor()
	opts.Diff = func(a, b interface{}) string {
		return " diff"
	}
	asDiff := got.NewWith(m, opts)
	asDiff.Eq("a", "b")
	m.check(`"a" ⦗not ==⦘ "b" diff`)

	{
		count := as.Count(2)
		count()
		m.cleanup()
		m.check(`Should count 2 times, but got 1`)

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
		m.check(`Should count 1 times, but got 2`)
	}
}
