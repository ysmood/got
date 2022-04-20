package diff_test

import (
	"strings"
	"testing"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/diff"
	"github.com/ysmood/got/lib/gop"
)

var setup = got.Setup(func(g got.G) {
	g.ErrorHandler = got.NewDefaultAssertionError(nil, nil)
})

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
	ts := diff.TokenizeText("a", "b")

	df := diff.Format(ts, diff.ThemeNone)

	g.Eq(df, ""+
		"1   - a\n"+
		"  1 + b\n"+
		"")
}

func TestTwoLines(t *testing.T) {
	g := setup(t)

	x, y := diff.TokenizeLine("abcdfghijklmn", "xxbcdfghixklmn")

	format := func(ts []*diff.Token) string {
		out := ""
		for _, t := range ts {
			switch t.Type {
			case diff.DelWords:
				out += "-" + t.Literal
			case diff.AddWords:
				out += "+" + t.Literal
			default:
				out += "." + t.Literal
			}
		}
		return out
	}

	g.Eq(format(x), "-a.b.c.d.f.g.h.i-j.k.l.m.n")
	g.Eq(format(y), "+x+x.b.c.d.f.g.h.i+x.k.l.m.n")
}

func TestColor(t *testing.T) {
	g := setup(t)

	out := diff.Diff("abc", "axc")

	g.Eq(gop.VisualizeANSI(out), `<45>@@ diff chunk @@<49>
<41>1   -<49> a<41>b<49>c
<42>  1 +<49> a<42>x<49>c

`)
}
