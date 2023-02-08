package diff_test

import (
	"context"
	"strings"
	"testing"

	"github.com/ysmood/gop"
	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/diff"
	"github.com/ysmood/got/lib/lcs"
)

var setup = got.Setup(func(g got.G) {
	g.ErrorHandler = got.NewDefaultAssertionError(nil, nil)
})

func split(s string) []string { return strings.Split(s, "") }

func TestDiff(t *testing.T) {
	g := setup(t)

	out := gop.StripANSI(diff.Diff("abc", "axc"))

	g.Eq(out, `@@ diff chunk @@
1   - abc
  1 + axc

`)
}

func TestFormat(t *testing.T) {
	g := setup(t)
	ts := diff.TokenizeText(
		g.Context(),
		strings.ReplaceAll("a b c d f g h h j q z", " ", "\n"),
		strings.ReplaceAll("a b c d e f g i j k r x y z", " ", "\n"),
	)

	df := diff.Format(ts, diff.ThemeNone)

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
		g.Context(),
		strings.ReplaceAll("a b c d f g h i j k l m n", " ", "\n"),
		strings.ReplaceAll("x b c d f g h i x k l m n", " ", "\n"),
	)

	lines := diff.ParseTokenLines(ts)
	lines = diff.Narrow(1, lines)
	ts = diff.SpreadTokenLines(lines)

	df := diff.Format(ts, diff.ThemeNone)

	g.Eq(df, ""+
		"@@ diff chunk @@\n"+
		"01    - a\n"+
		"   01 + x\n"+
		"02 02   b\n"+
		"\n"+
		"@@ diff chunk @@\n"+
		"08 08   i\n"+
		"09    - j\n"+
		"   09 + x\n"+
		"10 10   k\n"+
		"\n"+
		"")
}

func TestChunks0(t *testing.T) {
	g := setup(t)
	ts := diff.TokenizeText(
		g.Context(),
		strings.ReplaceAll("a b c", " ", "\n"),
		strings.ReplaceAll("a x c", " ", "\n"),
	)

	lines := diff.ParseTokenLines(ts)
	lines = diff.Narrow(-1, lines)
	ts = diff.SpreadTokenLines(lines)

	df := diff.Format(ts, diff.ThemeNone)

	g.Eq(df, ""+
		"@@ diff chunk @@\n"+
		"2   - b\n"+
		"  2 + x\n"+
		"\n"+
		"")
}

func TestNoDifference(t *testing.T) {
	g := setup(t)
	ts := diff.TokenizeText(g.Context(), "a", "b")

	df := diff.Format(ts, diff.ThemeNone)

	g.Eq(df, ""+
		"1   - a\n"+
		"  1 + b\n"+
		"")
}

func TestTwoLines(t *testing.T) {
	g := setup(t)

	format := func(ts []*diff.Token) string {
		out := ""
		for _, t := range ts {
			txt := strings.TrimSpace(strings.ReplaceAll(t.Literal, "", " "))
			switch t.Type {
			case diff.DelWords:
				out += "-" + txt
			case diff.AddWords:
				out += "+" + txt
			default:
				out += "=" + txt
			}
		}
		return out
	}

	check := func(x, y, ex, ey string) {
		t.Helper()

		tx, ty := diff.TokenizeLine(g.Context(),
			strings.ReplaceAll(x, " ", ""),
			strings.ReplaceAll(y, " ", ""))
		dx, dy := format(tx), format(ty)

		if dx != ex || dy != ey {
			t.Error("\n", dx, "\n", dy, "\n!=\n", ex, "\n", ey)
		}
	}

	check(
		" a b c d f g h i j k l m n",
		" x x b c d f g h i x k l m n",
		"-a=b c d f g h i-j=k l m n",
		"+x x=b c d f g h i+x=k l m n",
	)

	check(
		" 4 9 0 4 5 0 8 8 5 3",
		" 4 9 0 5 4 3 7 5 2",
		"=4 9 0 4 5-0 8 8 5 3",
		"=4 9 0+5=4+3 7=5+2",
	)

	check(
		" 4 9 0 4 5 0 8",
		" 4 9 0 5 4 3 7",
		"=4 9 0 4-5 0 8",
		"=4 9 0+5=4+3 7",
	)
}

func TestColor(t *testing.T) {
	g := setup(t)

	out := diff.Diff("abc", "axc")

	g.Eq(gop.VisualizeANSI(out), `<45><30>@@ diff chunk @@<39><49>
<31>1   -<39> a<31>b<39>c
<32>  1 +<39> a<32>x<39>c

`)
}

func TestCustomSplit(t *testing.T) {
	g := setup(t)

	ctx := context.WithValue(g.Context(), lcs.SplitKey, split)

	g.Eq(diff.TokenizeLine(ctx, "abc", "abc"))
}
