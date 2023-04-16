package lcs_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"runtime/debug"
	"testing"
	"time"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/lcs"
)

var setup = got.Setup(func(g got.G) {
	g.ErrorHandler = got.NewDefaultAssertionError(nil, nil)
})

func TestLCS(t *testing.T) {
	g := got.T(t)

	check := func(i int, x, y string) {
		t.Helper()

		e := func(msg ...interface{}) {
			t.Helper()
			t.Log(i, x, y)
			t.Error(msg...)
			t.FailNow()
		}

		defer func() {
			err := recover()
			if err != nil {
				debug.PrintStack()
				e(err)
			}
		}()

		xs, ys := lcs.NewChars(x), lcs.NewChars(y)

		s := xs.Sub(xs.YadLCS(context.Background(), ys))
		out := s.String()
		expected := lcs.StandardLCS(xs, ys).String()

		if !s.IsSubsequenceOf(xs) {
			e(s.String(), "is not subsequence of", x)
		}

		if !s.IsSubsequenceOf(ys) {
			e(s.String(), "is not subsequence of", y)
		}

		if len(out) != len(expected) {
			e("length of", out, "doesn't equal length of", expected)
		}
	}

	randStr := func() string {
		const c = 8
		b := make([]byte, c)
		for i := 0; i < c; i++ {
			b[i] = byte('a' + g.RandInt(0, c))
		}
		return string(b)
	}

	check(0, "", "")
	check(0, "", "a")

	for i := 1; i < 1000; i++ {
		check(i, randStr(), randStr())
	}
}

func TestLCSLongContentSmallChange(t *testing.T) {
	eq := func(x, y, expected string) {
		t.Helper()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		xs, ys := lcs.NewChars(x), lcs.NewChars(y)
		lcs := xs.YadLCS(ctx, ys)
		out := xs.Sub(lcs).String()

		if out != expected {
			t.Error(out, "!=", expected)
		}
	}

	x := bytes.Repeat([]byte("x"), 100000)
	y := bytes.Repeat([]byte("y"), 100000)
	eq(string(x), string(y), "")

	x[len(x)/2] = byte('a')
	y[len(y)/2] = byte('a')
	eq(string(x), string(y), "a")

	x[len(x)/2] = byte('y')
	y[len(y)/2] = byte('x')
	eq(string(x), string(y), "xy")
}

func TestContext(t *testing.T) {
	g := got.T(t)

	c := g.Context()
	c.Cancel()
	l := lcs.NewChars("abc").YadLCS(c, lcs.NewChars("abc"))
	g.Len(l, 0)
}

func TestLongRandom(_ *testing.T) {
	size := 10000
	x := randStr(size)
	y := randStr(size)

	c := context.Background()

	xs := lcs.NewChars(x)
	ys := lcs.NewChars(y)
	xs.YadLCS(c, ys)
}

func randStr(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return string(b)
}
