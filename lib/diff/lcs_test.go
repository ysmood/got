package diff_test

import (
	"strings"
	"testing"

	"github.com/ysmood/got/lib/diff"
)

func TestLCSString(t *testing.T) {
	g := setup(t)
	eq := func(x, y, expected string) {
		t.Helper()

		lcs := diff.LCS(diff.NewString(x), diff.NewString(y))
		g.Eq(diff.String(lcs).String(), expected)
	}

	eq("", "", "")
	eq("abc", "acb", "ab")
	eq("abc", "acbc", "abc")
	eq("abc", "xxx", "")
	eq("ac", "bc", "c")
	eq("gac", "agcat", "ga")
	eq("agcat", "gac", "ac")
}

func TestText(t *testing.T) {
	g := setup(t)
	g.Len(diff.NewText("a"), 1)
	g.Len(diff.NewText("a\n"), 2)
	g.Len(diff.NewText("a\n\n"), 3)
	g.Len(diff.NewText("\na"), 2)
}

func TestLCSText(t *testing.T) {
	g := setup(t)
	eq := func(x, y, expected string) {
		t.Helper()

		x = strings.Join(strings.Split(x, ""), "\n")
		y = strings.Join(strings.Split(y, ""), "\n")
		expected = strings.Join(strings.Split(expected, ""), "\n")

		lcs := diff.LCS(diff.NewText(x), diff.NewText(y))
		g.Eq(diff.Text(lcs).String(), expected)
	}

	eq("", "", "")
	eq("abc", "acb", "ab")
	eq("abc", "acbc", "abc")
	eq("abc", "xxx", "")
}
