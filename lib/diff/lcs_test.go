package diff_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/diff"
)

func TestReduce(t *testing.T) {
	eq := func(x, y string, e string) {
		t.Helper()

		out := diff.NewWords(diff.Split(x)).Reduce(diff.NewWords(diff.Split(y))).String()
		if out != e {
			t.Error(out, "!=", e)
		}
	}

	eq("", "", "")
	eq("", "a", "")
	eq("a", "", "")
	eq("abc", "abc", "abc")
	eq("abc", "acb", "abc")
	eq("abc", "acbc", "abc")
	eq("abc", "xxx", "")
	eq("ac", "bc", "c")
}

func TestCommon(t *testing.T) {
	eq := func(x, y string, el, er int) {
		t.Helper()

		l, r := diff.NewWords(diff.Split(x)).Common(diff.NewWords(diff.Split(y)))

		if l != el || r != er {
			t.Error(l, r, "!=", el, er)
		}
	}

	eq("", "", 0, 0)
	eq("", "a", 0, 0)
	eq("a", "", 0, 0)
	eq("abc", "abc", 3, 0)
	eq("abc", "acb", 1, 0)
	eq("abc", "acbc", 1, 2)
	eq("abc", "xxx", 0, 0)
	eq("ac", "bc", 0, 1)
}

func TestLCSString(t *testing.T) {
	eq := func(x, y, expected string) {
		t.Helper()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		lcs := diff.NewWords(split(x)).LCS(ctx, diff.NewWords(split(y)))
		out := lcs.String()

		if out != expected {
			t.Error(out, "!=", expected)
		}
	}

	eq("", "", "")
	eq("abc", "acb", "ac")
	eq("abc", "acbc", "abc")
	eq("abc", "xxx", "")
	eq("ac", "bc", "c")
	eq("gac", "agcat", "gc")
	eq("agcat", "gac", "ga")

	x := bytes.Repeat([]byte("x"), 100000)
	y := bytes.Repeat([]byte("y"), 100000)
	eq(string(x), string(y), "")

	x[len(x)/2] = byte('a')
	y[len(y)/2] = byte('a')
	eq(string(x), string(y), "a")

	x[len(x)/2] = byte('y')
	y[len(y)/2] = byte('x')
	eq(string(x), string(y), "yx")
}

func TestLCSRandomString(t *testing.T) {
	g := got.T(t)

	x := g.ReadFile("fixtures/rand_x.txt").String()
	y := g.ReadFile("fixtures/rand_y.txt").String()

	res := diff.NewWords(split(x)).LCS(g.Context(), diff.NewWords(split(y)))
	g.Eq(
		res.String(),
		g.ReadFile("fixtures/rand_lcs.txt").String(),
	)
}

func TestText(t *testing.T) {
	g := setup(t)
	g.Len(diff.NewLines("a"), 1)
	g.Len(diff.NewLines("a\n"), 2)
	g.Len(diff.NewLines("a\n\n"), 3)
	g.Len(diff.NewLines("\na"), 2)
}

func TestLCSText(t *testing.T) {
	g := setup(t)
	eq := func(x, y, expected string) {
		t.Helper()

		x = strings.Join(strings.Split(x, ""), "\n")
		y = strings.Join(strings.Split(y, ""), "\n")
		expected = strings.Join(strings.Split(expected, ""), "\n")

		lcs := diff.NewLines(x).LCS(context.Background(), diff.NewLines(y))
		g.Eq(lcs.String(), expected)
	}

	eq("", "", "")
	eq("abc", "acb", "ac")
	eq("abc", "acbc", "abc")
	eq("abc", "xxx", "")
}

func TestIsIn(t *testing.T) {
	g := got.T(t)

	y := diff.NewWords(split("abc"))

	g.True(diff.NewWords(split("ab")).IsSubsequenceOf(y))
	g.True(diff.NewWords(split("ac")).IsSubsequenceOf(y))
	g.True(diff.NewWords(split("bc")).IsSubsequenceOf(y))
	g.False(diff.NewWords(split("cb")).IsSubsequenceOf(y))
	g.False(diff.NewWords(split("ba")).IsSubsequenceOf(y))
	g.False(diff.NewWords(split("ca")).IsSubsequenceOf(y))
}

func TestBTreeFindGreater(t *testing.T) {
	g := got.T(t)

	check := func(x int, ey int, ef bool) {
		g.Helper()

		y, f := diff.BTreeFindGreater([]int{1, 3, 5, 7, 9, 11}, x)
		g.Eq(y, ey)
		g.Eq(f, ef)
	}

	check(4, 5, true)
	check(5, 7, true)
	check(1, 3, true)
	check(8, 9, true)
	check(11, 0, false)
}
