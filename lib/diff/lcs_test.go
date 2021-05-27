package diff_test

import (
	"strings"

	"github.com/ysmood/got/lib/diff"
)

func (t T) LCSString() {
	eq := func(x, y, expected string) {
		t.Helper()

		lcs := diff.LCS(diff.NewString(x), diff.NewString(y))
		t.Eq(diff.String(lcs).String(), expected)
	}

	eq("", "", "")
	eq("abc", "acb", "ab")
	eq("abc", "acbc", "abc")
	eq("abc", "xxx", "")
	eq("ac", "bc", "c")
	eq("gac", "agcat", "ga")
	eq("agcat", "gac", "ac")
}

func (t T) Text() {
	t.Len(diff.NewText("a"), 1)
	t.Len(diff.NewText("a\n"), 2)
	t.Len(diff.NewText("a\n\n"), 3)
	t.Len(diff.NewText("\na"), 2)
}

func (t T) LCSText() {
	eq := func(x, y, expected string) {
		t.Helper()

		x = strings.Join(strings.Split(x, ""), "\n")
		y = strings.Join(strings.Split(y, ""), "\n")
		expected = strings.Join(strings.Split(expected, ""), "\n")

		lcs := diff.LCS(diff.NewText(x), diff.NewText(y))
		t.Eq(diff.Text(lcs).String(), expected)
	}

	eq("", "", "")
	eq("abc", "acb", "ab")
	eq("abc", "acbc", "abc")
	eq("abc", "xxx", "")
}
