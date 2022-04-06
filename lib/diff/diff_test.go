package diff_test

import (
	"strings"
	"testing"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/diff"
)

var setup = got.Setup(func(g got.G) {
	g.ErrorHandler = got.NewDefaultAssertionError(false, false)
})

func TestFormat(t *testing.T) {
	g := setup(t)
	ts := diff.TokenizeText(
		strings.ReplaceAll("a b c d f g h h j q z", " ", "\n"),
		strings.ReplaceAll("a b c d e f g i j k r x y z", " ", "\n"),
	)

	df := diff.Format(ts, diff.NoTheme)

	g.Eq(df, ""+
		"01 01   a\n"+
		"02 02   b\n"+
		"03 03   c\n"+
		"04 04   d\n"+
		"   05 + e\n"+
		"05 06   f\n"+
		"06 07   g\n"+
		"07    - h\n"+
		"08    - h\n"+
		"   08 + i\n"+
		"09 09   j\n"+
		"10    - q\n"+
		"   10 + k\n"+
		"   11 + r\n"+
		"   12 + x\n"+
		"   13 + y\n"+
		"11 14   z\n"+
		"")
}

func TestDisconnectedChunks(t *testing.T) {
	g := setup(t)
	ts := diff.TokenizeText(
		strings.ReplaceAll("a b c d f g h i j k l m n", " ", "\n"),
		strings.ReplaceAll("x b c d f g h i x k l m n", " ", "\n"),
	)

	df := diff.Format(ts, diff.NoTheme)

	g.Eq(df, ""+
		"01    - a\n"+
		"   01 + x\n"+
		"02 02   b\n"+
		"03 03   c\n"+
		"04 04   d\n"+
		"05 05   f\n"+
		"06 06   g\n"+
		"07 07   h\n"+
		"08 08   i\n"+
		"09    - j\n"+
		"   09 + x\n"+
		"10 10   k\n"+
		"11 11   l\n"+
		"12 12   m\n"+
		"13 13   n\n"+
		"")
}

func TestNoDifference(t *testing.T) {
	g := setup(t)
	ts := diff.TokenizeText("a", "b")

	df := diff.Format(ts, diff.NoTheme)

	g.Eq(df, ""+
		"1   - a\n"+
		"  1 + b\n"+
		"")
}

func TestTwoLines(t *testing.T) {
	g := setup(t)
	x, y := diff.TokenizeLine("abc", "acx")

	g.Eq(x, []*diff.Token{
		{Type: diff.SameWords, Literal: "a"},
		{Type: diff.DelWords, Literal: "b"},
		{Type: diff.SameWords, Literal: "c"},
	})
	g.Eq(y, []*diff.Token{
		{Type: diff.SameWords, Literal: "a"},
		{Type: diff.SameWords, Literal: "c"},
		{Type: diff.AddWords, Literal: "x"},
	})
}

func TestColor(t *testing.T) {
	g := setup(t)
	g.Eq(diff.Diff("a", "b"), "1   \x1b[41m- \x1b[0m\x1b[41ma\x1b[0m\n  1 \x1b[42m+ \x1b[0m\x1b[42mb\x1b[0m\n")
}
