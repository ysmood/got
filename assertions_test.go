package got_test

import (
	"errors"
	"fmt"
	"os"
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

	as.Neq(1.1, 1)
	as.Neq([]int{1, 2}, []int{2, 1})

	as.Equal(1, 1)

	as.Lt(time.Millisecond, time.Second)
	as.Lte(1, 1)

	as.Gt(2, 1.5)
	as.Gte(2, 2.0)

	now := time.Now()
	as.Eq(now, now)
	as.Lt(now, now.Add(time.Second))
	as.Gt(now.Add(time.Second), now)

	as.True(true)
	as.False(false)

	as.Nil(nil)
	as.Nil((*int)(nil))
	as.Nil(os.Stat("go.mod"))
	as.NotNil(1)

	as.Regex(`\d\d`, "10")
	as.Has(`test`, "es")

	as.Len([]int{1, 2}, 2)

	as.Err(1, 2, errors.New("err"))
	as.Panic(func() { panic(1) })

	as.Is(1, 2)
	err := errors.New("err")
	as.Is(err, err)
	as.Is(fmt.Errorf("%w", err), err)

	as.Eq(1, 1).Must()
}

func TestAssertionErr(t *testing.T) {
	m := &mock{t: t}
	as := got.New(m)

	type data struct {
		A int
		S string
	}

	as.Eq(1, 2.0).Msg("not %s", "equal")
	m.check("not equal")

	as.Eq(data{1, "a"}, data{1, "b"})
	m.check("{1 a} <got_test.data> ⦗not ≂⦘ {1 b} <got_test.data>")

	as.Eq(true, "a")
	m.check(`true ⦗not ≂⦘ "a"`)

	as.Equal(1, 1.0)
	m.check("1 ⦗not ==⦘ 1 <float64>")

	as.Neq(1, 1)
	m.check("1 ⦗not ≠⦘ 1")

	as.Lt(1, 1)
	m.check("1 ⦗not <⦘ 1")
	as.Lte(2, 1)

	m.check("2 ⦗not ≤⦘ 1")
	as.Gt(1, 1)
	m.check("1 ⦗not >⦘ 1")
	as.Gte(1, 2)
	m.check("1 ⦗not ≥⦘ 2")

	as.True(false)
	m.check("⦗should be <true>⦘")
	as.False(true)
	m.check("⦗should be <false>⦘")

	as.Nil(1)
	m.check("⦗last value⦘ 1 ⦗should be <nil>⦘")
	as.Nil()
	m.check("⦗no args received⦘")
	as.NotNil(nil)
	m.check("⦗last value shouldn't be <nil>⦘")
	as.NotNil((*int)(nil))
	m.check("<*int> ⦗shouldn't be <nil>⦘")
	as.NotNil()
	m.check("⦗no args received⦘")

	as.Regex(`\d\d`, "aaa")
	m.check(`\d\d ⦗should match⦘ aaa`)
	as.Has(`test`, "x")
	m.check("test ⦗should has⦘ x")

	as.Len([]int{1, 2}, 3)
	m.check("⦗expect len⦘ 2 ⦗to be⦘ 3")

	as.Err(nil)
	m.check("⦗last value⦘ nil ⦗should be <error>⦘")
	as.Panic(func() {})
	m.check("⦗should panic⦘")
	as.Err()
	m.check("⦗no args received⦘")
	as.Err(1)
	m.check("⦗last value⦘ 1 ⦗should be <error>⦘")

	func() {
		defer func() {
			_ = recover()
		}()
		as.E(1, errors.New("E"))
	}()
	m.check("⦗last value⦘ E <*errors.errorString> ⦗should be <nil>⦘")

	as.Is(1, 2.2)
	m.check("1 ⦗should kind of⦘ 2.2 <float64>")
	as.Is(errors.New("a"), errors.New("b"))
	m.check("a <*errors.errorString> ⦗should in chain of⦘ b <*errors.errorString>")
}
