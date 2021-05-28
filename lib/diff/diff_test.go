package diff_test

import (
	"strings"
	"testing"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/diff"
)

type T struct {
	got.G
}

func Test(t *testing.T) {
	got.Each(t, func(t *testing.T) T {
		return T{got.NewWith(t, got.NoColor().NoDiff())}
	})
}

func (t T) Format() {
	ts := diff.TokenizeText(
		strings.ReplaceAll("a b c d f g h h j q z", " ", "\n"),
		strings.ReplaceAll("a b c d e f g i j k r x y z", " ", "\n"),
	)

	df := diff.Format(ts, diff.NoTheme)

	t.Eq(df, ""+
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

func (t T) DisconnectedChunks() {
	ts := diff.TokenizeText(
		strings.ReplaceAll("a b c d f g h i j k l m n", " ", "\n"),
		strings.ReplaceAll("x b c d f g h i x k l m n", " ", "\n"),
	)

	df := diff.Format(ts, diff.NoTheme)

	t.Eq(df, ""+
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

func (t T) NoDifference() {
	ts := diff.TokenizeText("a", "b")

	df := diff.Format(ts, diff.NoTheme)

	t.Eq(df, ""+
		"1   - a\n"+
		"  1 + b\n"+
		"")
}

func (t T) TwoLines() {
	x, y := diff.TokenizeLine("abc", "acx")

	t.Eq(x, []*diff.Token{
		{Type: diff.SameWords, Literal: "a"},
		{Type: diff.DelWords, Literal: "b"},
		{Type: diff.SameWords, Literal: "c"},
	})
	t.Eq(y, []*diff.Token{
		{Type: diff.SameWords, Literal: "a"},
		{Type: diff.SameWords, Literal: "c"},
		{Type: diff.AddWords, Literal: "x"},
	})
}

func (t T) Color() {
	t.Eq(diff.Diff("a", "b"), "1   \x1b[31m- \x1b[0m\x1b[31ma\n\x1b[0m  1 \x1b[32m+ \x1b[0m\x1b[32mb\n\x1b[0m")
}
